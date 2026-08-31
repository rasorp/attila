// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Get an Attila job registration run",
		Category:  "run",
		Args:      true,
		UsageText: "attila job register run get [options] [run-id]",
		Flags:     helper.ClientNamespaceFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					getCLIErrorMsg,
					fmt.Errorf("expected 1 argument, got %v", numArgs)), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			id, err := ulid.Parse(cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError(getCLIErrorMsg, err), 1)
			}

			getResp, _, err := client.JobRegisterRuns().Get(
				cliCtx.Context,
				id,
				&api.QueryOpts{Namespace: helper.NamespaceFromFlags(cliCtx)},
			)
			if err != nil {
				return cli.Exit(helper.FormatError(getCLIErrorMsg, err), 1)
			}

			outputRun(cliCtx, getResp.Run)
			return nil
		},
	}
}
