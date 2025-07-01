// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package register

import (
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/job/register/method"
	"github.com/rasorp/attila/internal/cmd/job/register/plan"
	"github.com/rasorp/attila/internal/cmd/job/register/rule"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "register",
		Usage:           "Control, plan, and execute Nomad job registrations",
		Category:        "job",
		HideHelpCommand: true,
		UsageText:       "attila job register <command> [options] [args]",
		Subcommands: []*cli.Command{
			method.Command(),
			plan.Command(),
			rule.Command(),
		},
	}
}
