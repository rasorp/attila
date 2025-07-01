// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package method

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/internal/helper/file"
	"github.com/rasorp/attila/pkg/api"
)

const (
	createErrorMsg = "failed to create Attila job registration method"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create an Attila job registration method",
		Category:  "method",
		Args:      true,
		UsageText: "attila job register method create [options] [method-spec]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(createErrorMsg, fmt.Errorf("expected 1 argument, got %v", numArgs)), 1)
			}

			var methodObj api.JobRegisterMethod

			if err := file.ParseConfig(cliCtx.Args().First(), &methodObj); err != nil {
				return cli.Exit(helper.FormatError(createErrorMsg, err), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			methodCreateResp, _, err := client.JobRegisterMethods().Create(context.Background(), &methodObj)
			if err != nil {
				return cli.Exit(helper.FormatError(createErrorMsg, err), 1)
			}

			outputMethod(cliCtx, methodCreateResp.Method)
			return nil
		},
	}
}
