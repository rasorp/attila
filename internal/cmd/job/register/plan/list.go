// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package plan

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List job registration plans",
		Category:  "plan",
		Args:      false,
		UsageText: "attila job register plan list [options]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			listResp, _, err := client.JobRegisterPlans().List(context.Background(), nil)
			if err != nil {
				return cli.Exit(helper.FormatError(listCLIErrorMsg, err), 1)
			}

			_, _ = fmt.Fprint(cliCtx.App.Writer, formatPlanList(listResp.Plans))
			_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")

			return nil
		},
	}
}

func formatPlanList(plans []*api.JobRegisterPlan) string {
	if len(plans) == 0 {
		return "No job registration plans found"
	}

	out := make([]string, 0, len(plans)+1)
	out = append(out, "ID|Job|Namespace")
	for _, plan := range plans {
		out = append(out, fmt.Sprintf(
			"%s|%s|%s",
			plan.ID, plan.JobID, plan.JobNamespace))
	}

	return helper.FormatList(out)
}
