// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package rule

import (
	"context"
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List Attila job registration rules",
		Category:  "rule",
		Args:      false,
		UsageText: "attila job register rule list [options]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			ruleListResp, _, err := client.JobRegisterRules().List(context.Background())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to list Attila job registration rules", err), 1)
			}

			_, _ = fmt.Fprint(cliCtx.App.Writer, formatRuleList(ruleListResp.Rules))
			_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")

			return nil
		},
	}
}

func formatRuleList(rules []*api.JobRegisterRuleStub) string {
	if len(rules) == 0 {
		return "No Attila job registration rules found"
	}

	out := make([]string, 0, len(rules)+1)
	out = append(out, "Name|Region Contexts")
	for _, rule := range rules {
		out = append(out, fmt.Sprintf(
			"%s|%v",
			rule.Name, formatRegionContexts(rule.RegionContexts)))
	}

	return helper.FormatList(out)
}

func formatRegionContexts(ctxs []api.JobRegisterRuleRegionContext) string {

	ctxStrings := make([]string, len(ctxs))

	for i, ctx := range ctxs {
		ctxStrings[i] = string(ctx)
	}

	return strings.Join(ctxStrings, ", ")
}
