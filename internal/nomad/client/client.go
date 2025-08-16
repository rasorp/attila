// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"sync"

	"github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/server/nomad"
)

type Clients struct {
	logger zerolog.Logger

	clients     map[string]*api.Client
	clientsLock sync.RWMutex
}

func New(logger zerolog.Logger) nomad.ClientController {
	return &Clients{
		logger:  logger.With().Str("component", "client_controller").Logger(),
		clients: make(map[string]*api.Client),
	}
}

func (c *Clients) RegionDelete(name string) {
	c.clientsLock.Lock()
	delete(c.clients, name)
	c.clientsLock.Unlock()
	c.logger.Debug().Str("region_name", name).Msg("deleted Nomad regional client")
}

func (c *Clients) RegionGet(name string) (*api.Client, error) {
	c.clientsLock.RLock()
	defer c.clientsLock.RUnlock()

	regionClient, ok := c.clients[name]
	if !ok {
		return nil, fmt.Errorf("no Nomad client found for region %q", name)
	}
	return regionClient, nil
}

func (c *Clients) RegionNum() int {
	c.clientsLock.RLock()
	defer c.clientsLock.RUnlock()

	return len(c.clients)
}

func (c *Clients) RegionSet(name string, client *api.Client) {
	c.clientsLock.Lock()
	c.clients[name] = client
	c.clientsLock.Unlock()

	c.logger.Debug().Str("region_name", name).Msg("created Nomad regional client")
}
