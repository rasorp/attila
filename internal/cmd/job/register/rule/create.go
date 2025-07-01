// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package rule

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/internal/helper/file"
	"github.com/rasorp/attila/pkg/api"
)

const (
	createErrorMsg = "failed to create Attila job registration rule"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create an Attila job registration rule",
		Category:  "rule",
		Args:      true,
		UsageText: "attila job register rule create [options] [rule-spec]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					createErrorMsg,
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			var ruleObj api.JobRegisterRule

			if err := file.ParseConfig(cliCtx.Args().First(), &ruleObj); err != nil {
				return cli.Exit(helper.FormatError(createErrorMsg, err), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			ruleCreateResp, _, err := client.JobRegisterRules().Create(context.Background(), &ruleObj)
			if err != nil {
				return cli.Exit(helper.FormatError(createErrorMsg, err), 1)
			}

			outputRule(cliCtx, ruleCreateResp.Rule)
			return nil
		},
	}
}
