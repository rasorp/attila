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

func getCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Detail an Attila region",
		Category:  "region",
		Args:      true,
		UsageText: "attila region get [options] [region-name]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					"failed to get Attila region",
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			regionResp, _, err := client.Regions().Get(context.Background(), cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to get Attila region", err), 1)
			}

			outputRegion(cliCtx, regionResp.Region)
			return nil
		},
	}
}
