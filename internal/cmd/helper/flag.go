// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package helper

import (
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/pkg/api"
)

const (
	addressCLIFlag   = "address"
	namespaceCLIFlag = "namespace"
)

func ClientFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Aliases: []string{"a"},
			Name:    addressCLIFlag,
			Value:   "http://127.0.0.1:8080",
			Usage:   "Attila server address to make API requests to",
		},
	}
}

func ClientNamespaceFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Aliases: []string{"a"},
			Name:    addressCLIFlag,
			Value:   "http://127.0.0.1:8080",
			Usage:   "Attila server address to make API requests to",
		},
		&cli.StringFlag{
			Aliases: []string{"n"},
			Name:    namespaceCLIFlag,
			Value:   "default",
			Usage:   "Attila namespace to make API requests to",
		},
	}
}

func ClientConfigFromFlags(ctx *cli.Context) *api.Config {

	defaultConfig := api.DefaultConfig()

	if addr := ctx.String(addressCLIFlag); addr != "" {
		defaultConfig.Address = addr
	}

	return defaultConfig
}

func NamespaceFromFlags(ctx *cli.Context) string { return ctx.String(namespaceCLIFlag) }
