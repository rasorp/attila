// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package plan

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Detail a job registration plan",
		Category:  "plan",
		Args:      true,
		UsageText: "attila job register plan get [options] [plan-id]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					getCLIErrorMsg,
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			id, err := ulid.Parse(cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError(getCLIErrorMsg, err), 1)
			}

			getReq := api.JobRegisterPlanGetReq{ID: id}

			getResp, _, err := client.JobRegisterPlans().Get(context.Background(), &getReq)
			if err != nil {
				return cli.Exit(helper.FormatError(getCLIErrorMsg, err), 1)
			}

			outputPlan(cliCtx, getResp.Plan)
			return nil
		},
	}
}
