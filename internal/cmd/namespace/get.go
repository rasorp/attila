// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Detail an Attila namespace",
		Category:  "namespace",
		Args:      true,
		UsageText: "attila namespace get [options] [namespace-name]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					"failed to get Attila namespace",
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			nsResp, _, err := client.Namespaces().Get(cliCtx.Context, cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to get Attila namespace", err), 1)
			}

			outputNamespace(cliCtx, nsResp.Namespace)
			return nil
		},
	}
}
