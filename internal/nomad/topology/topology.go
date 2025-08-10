// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package topology

import (
	"sync"

	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/nomad/client"
	"github.com/rasorp/attila/internal/server/nomad"
)

type Topology struct {
	clients *client.Clients
	logger  zerolog.Logger

	// regions stores the topology collector for the named region. Access should
	// use the lock for concurrent safety.
	regions     map[string]*region
	regionsLock sync.RWMutex
}

func New(logger *zerolog.Logger, clients *client.Clients) nomad.TopologyController {
	return &Topology{
		clients: clients,
		logger:  logger.With().Str("component", "topology").Logger(),
		regions: make(map[string]*region),
	}
}

func (c *Topology) RegionSet(name string, _ *api.Client) {
	c.regionsLock.Lock()
	defer c.regionsLock.Unlock()

	// If the region is currently tracked, it might be that the region state
	// specification has been modified. While this might change the Nomad API
	// client, it doesn't need to be propagated as the region collector pulls
	// from the client store on each collection.
	if _, ok := c.regions[name]; ok {
		return
	}

	c.regions[name] = newRegion(name, c.clients, c.logger)
	go c.regions[name].run()
}

func (c *Topology) RegionDelete(name string) {
	c.regionsLock.Lock()
	defer c.regionsLock.Unlock()

	if capacityRunner, ok := c.regions[name]; ok {
		capacityRunner.stop()
		delete(c.regions, name)
	}
}

func (c *Topology) GetTopologies() []*nomad.Overview {
	c.regionsLock.RLock()
	defer c.regionsLock.RUnlock()

	// Create our slice with the allocated length but an initial size of zero
	// because we are iterating a map, so we must append rather than inset into
	// specific indexes.
	out := make([]*nomad.Overview, 0, len(c.regions))

	// Iterate over the stored collectors. If they have not been able to collect
	// topology data, due to API errors, the overview will be nil. We want to
	// ensure this isn't added to the return object.
	for _, regionController := range c.regions {
		if overview := regionController.getOverviewResult(); overview != nil {
			out = append(out, overview)
		}
	}

	return out
}

func (c *Topology) GetTopology(name string) *nomad.Topology {
	c.regionsLock.RLock()
	defer c.regionsLock.RUnlock()

	if val, ok := c.regions[name]; ok {
		return val.getResult()
	}

	return nil
}
