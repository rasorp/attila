// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package rule

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "rule",
		Category:        "register",
		Usage:           "ister Attila job registration rules",
		HideHelpCommand: true,
		UsageText:       "attila job register rule <command> [options] [args]",
		Subcommands: []*cli.Command{
			createCommand(),
			deleteCommand(),
			getCommand(),
			listCommand(),
		},
	}
}

func outputRule(cliCtx *cli.Context, r *api.JobRegisterRule) {
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
		fmt.Sprintf("Name|%s", r.Name),
		fmt.Sprintf("Region Contexts|%s", contextsAsString(r.RegionContexts)),
		fmt.Sprintf("Region Pickers|%s", formatRegionPicker(r.RegionPickers)),
		fmt.Sprintf("Create Time|%s", helper.FormatTime(r.Metadata.CreateTime)),
		fmt.Sprintf("Update Time|%s", helper.FormatTime(r.Metadata.UpdateTime)),
	}))
	_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")
}

func contextsAsString(ctxs []*api.JobRegisterRuleRegionContext) string {
	var s []string
	for _, regionCtx := range ctxs {
		s = append(s, string(*regionCtx))
	}
	return strings.Join(s, ", ")
}

func formatRegionPicker(rp []*api.JobRegisterRegionPicker) string {
	if len(rp) == 0 {
		return ""
	}

	out := make([]string, 0, len(rp))
	for _, picker := range rp {
		if picker == nil || picker.Provider == "" {
			continue
		}
		out = append(out, fmt.Sprintf("%s::%s", picker.Provider, picker.Name))
	}

	return strings.Join(out, ", ")
}
