// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package topology

import (
	"sync"
	"time"

	"github.com/hashicorp/nomad/api"
	"go.uber.org/zap"

	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/server/nomad"
)

// defaultCollectionInterval is the base interval used for the collection
// ticker. While this is currently hardcoded, we probably should move this to a
// configurable parameter at some point.
var defaultCollectionInterval = 1 * time.Minute

type region struct {
	name    string
	clients *client.Clients
	logger  *zap.Logger

	// result stores the last fetched result of the region topology. All access
	// should use the lock, as the object is concurrently written/read via a
	// number of routines.
	result     *nomad.Topology
	resultLock sync.RWMutex

	// shutdownCh is used to instruct the long-lived routine to shut down.
	shutdownCh chan struct{}
}

func newRegion(name string, clients *client.Clients, logger *zap.Logger) *region {
	return &region{
		name:       name,
		clients:    clients,
		logger:     logger.With(zap.String("region", name)),
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

	r.logger.Info(
		"starting periodic collector",
		zap.Int64("interval_ms", defaultCollectionInterval.Milliseconds()),
	)

	ticker := time.NewTicker(defaultCollectionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.runExecute()
		case <-r.shutdownCh:
			r.logger.Info("shutting down topology collector")
			return
		}
	}
}

func (r *region) runExecute() {

	// Track the start time, so we can monitor how long it takes for the
	// collection to run.
	startTime := time.Now()
	r.logger.Info("performing execution of data collection")

	apiClient, err := r.clients.Get(r.name)
	if err != nil {
		r.logger.Error("failed to get API client", zap.Error(err))
		return
	}

	result := nomad.NewTopology(r.name)

	if err := r.executeAgentMembers(apiClient, result); err != nil {
		r.logger.Error("failed to process server topology", zap.Error(err))
		return
	}

	if err := r.executeNodes(apiClient, result); err != nil {
		r.logger.Error("failed to process node topology", zap.Error(err))
		return
	}

	r.resultLock.Lock()
	r.result = result
	r.resultLock.Unlock()

	r.logger.Info(
		"finished execution of data collection",
		zap.Int64("dur", int64(time.Since(startTime))),
	)
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
