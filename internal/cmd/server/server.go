// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"github.com/urfave/cli/v2"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "server",
		Usage:           "Run and control Attila servers",
		HideHelpCommand: true,
		UsageText:       "attila server <command> [options] [args]",
		Subcommands: []*cli.Command{
			runCommand(),
		},
	}
}
