// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func deleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete an Attila namespace",
		Category:  "namespace",
		Args:      true,
		UsageText: "attila namespace delete [options] [namespace-name]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					"failed to delete Attila namespace",
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			_, err := client.Namespaces().Delete(cliCtx.Context, cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to delete Attila namespace", err), 1)
			}

			_, _ = fmt.Fprintf(cliCtx.App.Writer, "successfully deleted Attila namespace %q", cliCtx.Args().First())
			return nil
		},
	}
}
