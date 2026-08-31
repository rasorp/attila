// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List Attila job registration runs",
		Category:  "run",
		UsageText: "attila job register run list [options]",
		Flags:     helper.ClientNamespaceFlags(),
		Action: func(cliCtx *cli.Context) error {
			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			listResp, _, err := client.JobRegisterRuns().List(
				cliCtx.Context,
				&api.QueryOpts{Namespace: helper.NamespaceFromFlags(cliCtx)},
			)
			if err != nil {
				return cli.Exit(helper.FormatError(listCLIErrorMsg, err), 1)
			}

			outputRuns(cliCtx, listResp.Runs)
			return nil
		},
	}
}
