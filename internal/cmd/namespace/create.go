// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/internal/helper/file"
	"github.com/rasorp/attila/pkg/api"
)

const (
	commandCLIErrorMsg = "failed to create Attila namespace"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create an Attila namespace",
		Category:  "namespace",
		Args:      true,
		UsageText: "attila namespace create [options] [namespace-spec]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs < 1 {
				return cli.Exit(helper.FormatError(
					commandCLIErrorMsg,
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			var namespace api.Namespace

			if err := file.ParseConfig(cliCtx.Args().First(), &namespace); err != nil {
				return cli.Exit(helper.FormatError(commandCLIErrorMsg, err), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			req := api.NamespaceCreateReq{Namespace: &namespace}

			nsResp, _, err := client.Namespaces().Create(cliCtx.Context, &req)
			if err != nil {
				return cli.Exit(helper.FormatError(commandCLIErrorMsg, err), 1)
			}

			outputNamespace(cliCtx, nsResp.Namespace)
			return nil
		},
	}
}
