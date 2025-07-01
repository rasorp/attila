// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package region

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func deleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete an Attila region",
		Category:  "region",
		Args:      true,
		UsageText: "attila region delete [options] [region-name]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					"failed to delete Attila region",
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			_, err := client.Regions().Delete(context.Background(), cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to delete Attila region", err), 1)
			}

			_, _ = fmt.Fprintf(cliCtx.App.Writer, "successfully deleted Attila region %q", cliCtx.Args().First())
			return nil
		},
	}
}
