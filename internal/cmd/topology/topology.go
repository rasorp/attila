// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package topology

import "github.com/urfave/cli/v2"

func Command() *cli.Command {
	return &cli.Command{
		Name:            "topology",
		Usage:           "View collected region topologies",
		HideHelpCommand: true,
		UsageText:       "attila topology <command> [options] [args]",
		Subcommands: []*cli.Command{
			getCommand(),
			listCommand(),
		},
	}
}
