// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package job

import (
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/job/register"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "job",
		Usage:           "Manipulate and control Nomad jobs",
		HideHelpCommand: true,
		UsageText:       "attila job <command> [options] [args]",
		Subcommands: []*cli.Command{
			register.Command(),
		},
	}
}
