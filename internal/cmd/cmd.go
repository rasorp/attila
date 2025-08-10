// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/internal/cmd/job"
	"github.com/rasorp/attila/internal/cmd/region"
	"github.com/rasorp/attila/internal/cmd/server"
	"github.com/rasorp/attila/internal/cmd/topology"
	"github.com/rasorp/attila/internal/version"
)

func main() {

	cli.VersionPrinter = func(cliCtx *cli.Context) {
		_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
			fmt.Sprintf("Version|%s", cliCtx.App.Version),
			fmt.Sprintf("Build Time|%s", version.BuildTime),
			fmt.Sprintf("Build Commit|%s", version.BuildCommit),
		}))
		_, _ = fmt.Fprint(cliCtx.App.Writer, "\n")
	}

	cliApp := cli.App{
		Commands: []*cli.Command{
			job.Command(),
			region.Command(),
			server.Command(),
			topology.Command(),
		},
		Name:  "attila",
		Usage: "Meta application and scheduling for HashiCorp Nomad",
		Description: `Attila is a meta application for the Nomad workload orchestrator. It provides
high level abstractions and operator tooling which should not exist in Nomad,
and in particular, aims to help operate non-federate, cell based cluster
deployments.`,
		Version:         version.Get(),
		HideHelpCommand: true,
	}

	// The CLI application should handle incorrect flags and arguments and
	// output-friendly responses. Printing this error will duplicate the CLI
	// response, so it is safe to throw away.
	_ = cliApp.Run(os.Args)
}
