// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package method

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
		Usage:     "Detail an Attila job registration method",
		Category:  "method",
		Args:      true,
		UsageText: "attila job register method get [options] [method-name]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					"failed to get Attila job registration method",
					fmt.Errorf("expected 1 argument, got %v", numArgs)), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			methodResp, _, err := client.JobRegisterMethods().Get(context.Background(), cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to get Attila job registration method", err), 1)
			}

			outputMethod(cliCtx, methodResp.Method)
			return nil
		},
	}
}
