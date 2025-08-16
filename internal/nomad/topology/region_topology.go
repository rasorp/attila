// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package topology

import (
	"sync"
	"time"

	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/server/nomad"
)

// defaultCollectionInterval is the base interval used for the collection
// ticker. While this is currently hardcoded, we probably should move this to a
// configurable parameter at some point.
var defaultCollectionInterval = 1 * time.Minute

type region struct {
	name    string
	clients nomad.ClientController
	logger  zerolog.Logger

	// result stores the last fetched result of the region topology. All access
	// should use the lock, as the object is concurrently written/read via a
	// number of routines.
	result     *nomad.Topology
	resultLock sync.RWMutex

	// shutdownCh is used to instruct the long-lived routine to shut down.
	shutdownCh chan struct{}
}

func newRegion(name string, clients nomad.ClientController, logger zerolog.Logger) *region {
	return &region{
		name:       name,
		clients:    clients,
		logger:     logger.With().Str("region", name).Logger(),
		shutdownCh: make(chan struct{}),
	}
}

func (r *region) getResult() *nomad.Topology {
	r.resultLock.RLock()
	defer r.resultLock.RUnlock()
	return r.result
}

func (r *region) getOverviewResult() *nomad.Overview {
	r.resultLock.RLock()
	defer r.resultLock.RUnlock()

	if r.result != nil {
		return r.result.Overview
	}

	return nil
}

func (r *region) run() {

	// Perform an initial collection as soon as the region topology collector is
	// created. This means we do not have to wait for the ticker to fire before
	// we populate the result.
	r.runExecute()

	r.logger.Info().
		Dur("interval", defaultCollectionInterval).
		Msg("starting periodic collector")

	ticker := time.NewTicker(defaultCollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.runExecute()
		case <-r.shutdownCh:
			r.logger.Info().Msg("shutting down topology collector")
			return
		}
	}
}

func (r *region) runExecute() {

	// Track the start time, so we can monitor how long it takes for the
	// collection to run.
	startTime := time.Now()
	r.logger.Info().Msg("performing execution of data collection")

	apiClient, err := r.clients.RegionGet(r.name)
	if err != nil {
		r.logger.Error().Err(err).Msg("failed to get API client")
		return
	}

	result := nomad.NewTopology(r.name)

	if err := r.executeAgentMembers(apiClient, result); err != nil {
		r.logger.Error().Err(err).Msg("failed to process server topology")
		return
	}

	if err := r.executeNodes(apiClient, result); err != nil {
		r.logger.Error().Err(err).Msg("failed to process node topology")
		return
	}

	r.resultLock.Lock()
	r.result = result
	r.resultLock.Unlock()

	r.logger.Info().
		TimeDiff("dur", time.Now(), startTime).
		Msg("finished execution of data collection")
}

func (r *region) executeAgentMembers(client *api.Client, result *nomad.Topology) error {

	members, err := client.Agent().Members()
	if err != nil {
		return err
	}

	for _, member := range members.Members {
		result.AddServer(member)
	}

	return nil
}

func (r *region) executeNodes(client *api.Client, result *nomad.Topology) error {

	nodeList, _, err := client.Nodes().List(&api.QueryOptions{
		Params: map[string]string{"resources": "true"},
	})
	if err != nil {
		return err
	}

	for _, node := range nodeList {

		nodeAllocs, _, err := client.Nodes().Allocations(node.ID, nil)
		if err != nil {
			return err
		}

		result.AddNode(node, nodeAllocs)
	}

	return nil
}

func (r *region) stop() { close(r.shutdownCh) }
