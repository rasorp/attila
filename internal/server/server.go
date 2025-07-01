// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"github.com/rasorp/attila/internal/helper/pointer"
	"github.com/rasorp/attila/internal/logger"
	nomadControler "github.com/rasorp/attila/internal/nomad"
	serverHTTP "github.com/rasorp/attila/internal/server/http"
	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/server/state"
	stateBackend "github.com/rasorp/attila/internal/state"
)

type Server struct {
	baseLogger   *zerolog.Logger
	serverLogger *zerolog.Logger
	srvs         []*httpServer
	state        state.State

	// nomadController
	nomadController nomad.Controller
}

type httpServer struct {
	logger *zerolog.Logger
	ln     net.Listener
	mux    *chi.Mux
	server *http.Server
}

func NewServer(cfg *Config) (*Server, error) {

	zerologger, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	stateBackend, err := stateBackend.NewBackend(cfg.State)
	if err != nil {
		return nil, fmt.Errorf("failed to setup state: %w", err)
	}

	server := Server{
		baseLogger:      zerologger,
		serverLogger:    pointer.Of(zerologger.With().Str("component", "server").Logger()),
		state:           stateBackend,
		nomadController: nomadControler.NewController(zerologger),
	}

	server.serverLogger.Info().Msg("successfully setup state backend")

	for _, bind := range cfg.HTTP.Binds {

		serverLogger := server.serverLogger.With().
			Str("address", bind.Addr).Logger()

		srv := httpServer{
			logger: &serverLogger,
			mux:    serverHTTP.NewRouter(serverLogger, cfg.HTTP.AccessLogLevel, stateBackend, server.nomadController),
		}

		// Configure the HTTP server to the most basic level.
		srv.server = &http.Server{
			Addr:         bind.Addr,
			Handler:      srv.mux,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 15 * time.Second,
			IdleTimeout:  15 * time.Second,
		}

		parsedURL, err := url.Parse(srv.server.Addr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse bind address: %w", err)
		}

		network := "tcp"
		if parsedURL.Scheme == "unix" {
			network = parsedURL.Scheme
		}

		ln, err := net.Listen(network, parsedURL.Host)
		if err != nil {
			return nil, fmt.Errorf("failed to setup HTTP listener: %w", err)
		}
		srv.ln = ln

		server.srvs = append(server.srvs, &srv)
		serverLogger.Info().Msg("successfully setup HTTP server")
	}

	return &server, nil
}

// Run is used to serve the HTTP server. The function will block and should be
// run via a go-routine. Unless http.Server.Serve panics/fails, the server can
// be stopped by calling the Stop function.
func (s *Server) Start() {
	for _, srv := range s.srvs {
		srv.logger.Info().Msg("server now listening for connections")
		go func() {
			_ = srv.server.Serve(srv.ln)
		}()
	}
}

func (s *Server) Stop() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for _, srv := range s.srvs {
		if err := srv.server.Shutdown(ctx); err != nil {
			srv.logger.Error().Err(err).Msg("failed to gracefully shutdown HTTP server")
		} else {
			srv.logger.Info().Msg("successfully shutdown HTTP server")
		}
	}
}

func (s *Server) WaitForSignals() {

	signalCh := make(chan os.Signal, 3)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Wait to receive a signal. This blocks until we are notified.
	for {
		s.serverLogger.Debug().Msg("wait for signal handler started")

		sig := <-signalCh
		s.serverLogger.Info().Str("signal", sig.String()).Msg("received signal")

		// Check the signal we received. If it was a SIGHUP when the
		// functionality is added, we perform the reload tasks and then
		// continue to wait for another signal. Everything else means exit.
		switch sig {
		case syscall.SIGHUP:
		default:
			s.Stop()
			return
		}
	}
}
