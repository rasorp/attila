// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/internal/helper/file"
	"github.com/rasorp/attila/internal/server"
)

func runCommand() *cli.Command {
	return &cli.Command{
		Name:     "run",
		Usage:    "Run an Attila server",
		Category: "server",
		Flags:    runFlags(),
		Action: func(cliCtx *cli.Context) error {

			cfg, err := generateRunConfig(cliCtx)
			if err != nil {
				return cli.Exit(helper.FormatError("failed to run an Attila server", err), 1)
			}

			srv, err := server.NewServer(cfg)
			if err != nil {
				return cli.Exit(helper.FormatError("failed to run an Attila server", err), 1)
			}
			srv.Start()
			srv.WaitForSignals()
			return nil
		},
	}
}

func runFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "config",
			Value: cli.NewStringSlice(),
			Usage: "The path to a config file",
		},
		&cli.StringFlag{
			Name:  "http-access-log-level",
			Value: "info",
			Usage: "The log verbosity to use for HTTP access logs",
		},
		&cli.StringSliceFlag{
			Name:  "http-bind-address",
			Value: cli.NewStringSlice(),
			Usage: "The HTTP/HTTPS/UNIX bind address for the server to use",
		},
		&cli.StringFlag{
			Name:  "log-level",
			Value: "info",
			Usage: "The log verbosity to use",
		},
		&cli.StringFlag{
			Name:  "log-format",
			Value: "json",
			Usage: "The format of log entires",
		},
		&cli.BoolFlag{
			Name:  "log-colour",
			Value: false,
			Usage: "Colourize logging output",
		},
		&cli.BoolFlag{
			Name:  "log-include-line",
			Value: false,
			Usage: "Add file:line of the caller to each log entry",
		},
		&cli.BoolFlag{
			Name:  "state-memory-enabled",
			Value: true,
			Usage: "Enable the memory state backend",
		},
	}
}

func generateRunConfig(cliCtx *cli.Context) (*server.Config, error) {

	defaultCfg := server.DefaultConfig()

	//
	for _, configFile := range cliCtx.StringSlice("config") {

		var serverCfg server.Config

		if err := file.ParseConfig(configFile, serverCfg); err != nil {
			return nil, err
		}

		defaultCfg = defaultCfg.Merge(&serverCfg)
	}

	//
	if lvl := cliCtx.String("http-access-log-level"); lvl != "" {
		defaultCfg.HTTP.AccessLogLevel = lvl
	}
	if len(cliCtx.StringSlice("http-bind-address")) > 0 {
		defaultCfg.HTTP.Binds = make([]*server.BindConfig, len(cliCtx.StringSlice("http-bind-address")))
		for i, addr := range cliCtx.StringSlice("http-bind-address") {
			defaultCfg.HTTP.Binds[i] = &server.BindConfig{Addr: addr}
		}
	}

	//
	if lvl := cliCtx.String("log-level"); lvl != "" {
		defaultCfg.Log.Level = lvl
	}
	if format := cliCtx.String("log-format"); format != "" {
		defaultCfg.Log.Format = format
	}
	if colour := cliCtx.Bool("log-colour"); colour {
		defaultCfg.Log.Colour = &colour
	}
	if line := cliCtx.Bool("log-include-line"); line {
		defaultCfg.Log.IncludeLine = &line
	}

	//
	if memoryState := cliCtx.Bool("state-memory-enabled"); memoryState {
		defaultCfg.State.Memory.Enabled = &memoryState
	}

	if err := defaultCfg.Validate(); err != nil {
		return nil, err
	}

	return defaultCfg, nil
}
