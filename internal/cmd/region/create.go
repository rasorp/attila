// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package region

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/internal/helper/file"
	"github.com/rasorp/attila/pkg/api"
)

const (
	commandCLIErrorMsg = "failed to create Attila region"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create an Attila region",
		Category:  "region",
		Args:      true,
		UsageText: "attila region create [options] [region-spec]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					commandCLIErrorMsg,
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			var region api.Region

			if err := file.ParseConfig(cliCtx.Args().First(), &region); err != nil {
				return cli.Exit(helper.FormatError(commandCLIErrorMsg, err), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			req := api.RegionCreateReq{Region: &region}

			regionCreateResp, _, err := client.Regions().Create(context.Background(), &req)
			if err != nil {
				return cli.Exit(helper.FormatError(commandCLIErrorMsg, err), 1)
			}

			outputRegion(cliCtx, regionCreateResp.Region)
			return nil
		},
	}
}
