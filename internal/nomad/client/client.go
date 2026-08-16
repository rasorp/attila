// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"fmt"
	"sync"

	"github.com/hashicorp/nomad/api"
	"go.uber.org/zap"
)

type Clients struct {
	logger *zap.Logger

	clients     map[string]*api.Client
	clientsLock sync.RWMutex
}

func New(logger *zap.Logger) *Clients {
	return &Clients{
		logger:  logger.Named("nomad_client"),
		clients: make(map[string]*api.Client),
	}
}

func (c *Clients) Delete(name string) {
	c.clientsLock.Lock()
	delete(c.clients, name)
	c.clientsLock.Unlock()
	c.logger.Debug("deleted Nomad regional client", zap.String("region_name", name))
}

func (c *Clients) Get(name string) (*api.Client, error) {
	c.clientsLock.RLock()
	defer c.clientsLock.RUnlock()

	regionClient, ok := c.clients[name]
	if !ok {
		return nil, fmt.Errorf("no Nomad client found for region %q", name)
	}
	return regionClient, nil
}

func (c *Clients) Num() int {
	c.clientsLock.RLock()
	defer c.clientsLock.RUnlock()

	return len(c.clients)
}

func (c *Clients) Set(name string, client *api.Client) {
	c.clientsLock.Lock()
	c.clients[name] = client
	c.clientsLock.Unlock()

	c.logger.Debug("created Nomad regional client", zap.String("region_name", name))
}
