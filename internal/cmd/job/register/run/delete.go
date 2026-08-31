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

func deleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete an Attila job registration run",
		Category:  "run",
		Args:      true,
		UsageText: "attila job register run delete [options] [run-id]",
		Flags:     helper.ClientNamespaceFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					deleteCLIErrorMsg,
					fmt.Errorf("expected 1 argument, got %v", numArgs)), 1)
			}

			id, err := ulid.Parse(cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError(deleteCLIErrorMsg, err), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			_, err = client.JobRegisterRuns().Delete(
				cliCtx.Context,
				&api.JobsRegisterRunsDeleteReq{ID: id},
				&api.WriteOpts{Namespace: helper.NamespaceFromFlags(cliCtx)},
			)
			if err != nil {
				return cli.Exit(helper.FormatError(deleteCLIErrorMsg, err), 1)
			}

			_, _ = fmt.Fprintf(
				cliCtx.App.Writer,
				"successfully deleted Attila job registration run %q",
				cliCtx.Args().First(),
			)
			return nil
		},
	}
}
