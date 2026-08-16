// Copyright James Rasell 2025, 2026
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
	"go.uber.org/zap"

	"github.com/rasorp/attila/internal/logger"
	nomadControler "github.com/rasorp/attila/internal/nomad"
	serverHTTP "github.com/rasorp/attila/internal/server/http"
	"github.com/rasorp/attila/internal/server/nomad"
	"github.com/rasorp/attila/internal/store"
	storebackend "github.com/rasorp/attila/internal/store/backend"
)

type Server struct {
	baseLogger   *zap.Logger
	serverLogger *zap.Logger
	srvs         []*httpServer
	state        store.State

	// nomadController
	nomadController nomad.Controller
}

type httpServer struct {
	logger *zap.Logger
	ln     net.Listener
	mux    *chi.Mux
	server *http.Server
}

func NewServer(cfg *Config) (*Server, error) {

	baseLogger, err := logger.New(cfg.Log)
	if err != nil {
		return nil, fmt.Errorf("failed to setup logger: %w", err)
	}

	backend, err := storebackend.NewBackend(cfg.State)
	if err != nil {
		return nil, fmt.Errorf("failed to setup state: %w", err)
	}

	server := Server{
		baseLogger:      baseLogger,
		serverLogger:    baseLogger.Named("server"),
		state:           backend,
		nomadController: nomadControler.NewController(baseLogger),
	}

	server.serverLogger.Info("successfully setup state backend")

	if err := server.restore(); err != nil {
		return nil, fmt.Errorf("failed to perform server restore: %w", err)
	}

	for _, bind := range cfg.HTTP.Binds {

		serverLogger := server.serverLogger.With(
			zap.String("address", bind.Addr),
		)

		srv := httpServer{
			logger: serverLogger,
			mux:    serverHTTP.NewRouter(serverLogger, cfg.HTTP.AccessLogLevel, backend, server.nomadController),
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
		listenAddr := parsedURL.Host
		if parsedURL.Scheme == "unix" {
			network = parsedURL.Scheme
		}

		ln, err := net.Listen(network, listenAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to setup HTTP listener: %w", err)
		}
		srv.ln = ln

		server.srvs = append(server.srvs, &srv)
		serverLogger.Info("successfully setup HTTP server")
	}

	return &server, nil
}

// restore handles restoration of Attila systems once the state backend has been
// set up and is accessible.
func (s *Server) restore() error {

	// List all the regions within our state, so we can restore the API clients.
	regionList, err := s.state.Region().List(nil)
	if err != nil {
		return err
	}

	// Iterate the regions stored within the state and restore the controller
	// client. If we are unable to create the API client, we log the problem but
	// continue. Causing the server to exit here would require operators to
	// manually intervene to restore the server. This way, the server will be
	// able to start and the impacted region configuration fixed when possible.
	for _, region := range regionList.Regions {

		apiClient, err := region.GenerateNomadClient()
		if err != nil {
			s.serverLogger.Error(
				"failed to restore region client",
				zap.String("region_name", region.Name),
				zap.Error(err),
			)
			continue
		}

		s.nomadController.RegionSet(region.Name, apiClient)
	}

	return nil
}

// Start is used to serve the HTTP server. The function will block and should be
// run via a go-routine. Unless http.Server.Serve panics/fails, the server can
// be stopped by calling the Stop function.
func (s *Server) Start() {
	for _, srv := range s.srvs {
		srv.logger.Info("server now listening for connections")
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
			srv.logger.Error("failed to gracefully shutdown HTTP server", zap.Error(err))
		} else {
			srv.logger.Info("successfully shutdown HTTP server")
		}

		_ = srv.ln.Close()
	}
}

func (s *Server) WaitForSignals() {

	signalCh := make(chan os.Signal, 3)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	// Wait to receive a signal. This blocks until we are notified.
	for {
		s.serverLogger.Debug("wait for signal handler started")

		sig := <-signalCh
		s.serverLogger.Info("received signal", zap.String("signal", sig.String()))

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
