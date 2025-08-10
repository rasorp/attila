// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/rasorp/attila/internal/helper/pointer"
	"github.com/rasorp/attila/internal/helper/test/mock"
	"github.com/rasorp/attila/internal/server/state"
	store "github.com/rasorp/attila/internal/state"
)

func TestServer_restore(t *testing.T) {

	cfg := DefaultConfig()

	cfg.State.File = &store.FileConfig{
		Enable: pointer.Of(true),
		Path:   t.TempDir(),
	}

	// Start an initial server and write a region to state. This acts as the
	// restore point for testing.
	startServer, err := NewServer(cfg)
	must.NoError(t, err)
	must.NotNil(t, startServer)

	_, err = startServer.state.Region().Create(
		&state.RegionCreateReq{Region: mock.Region()},
	)
	must.Nil(t, err)

	// Stop the server, so we can free the listener.
	startServer.Stop()

	// Build a new server and test that the Nomad controller has the expected
	// number of regions being tracked.
	restoreServer, err := NewServer(cfg)
	must.NoError(t, err)
	must.NotNil(t, restoreServer)

	must.Eq(t, 1, restoreServer.nomadController.RegionNum())

	restoreServer.Stop()
}
