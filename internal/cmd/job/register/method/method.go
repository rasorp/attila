// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package method

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "method",
		Category:        "register",
		Usage:           "Administer Attila job registration methods",
		HideHelpCommand: true,
		UsageText:       "attila job register method <command> [options] [args]",
		Subcommands: []*cli.Command{
			createCommand(),
			deleteCommand(),
			getCommand(),
			listCommand(),
		},
	}
}

func outputMethod(cliCtx *cli.Context, m *api.JobRegisterMethod) {
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
		fmt.Sprintf("Name|%s", m.Name),
		fmt.Sprintf("Selector|%s", m.Selector),
		fmt.Sprintf("Create Time|%s", helper.FormatTime(m.Metadata.CreateTime)),
		fmt.Sprintf("Update Time|%s", helper.FormatTime(m.Metadata.UpdateTime)),
	}))
	_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n\n")

	ruleList := make([]string, 0, len(m.Rules)+1)
	ruleList = append(ruleList, "Rules")
	for _, rule := range m.Rules {
		ruleList = append(ruleList, "  - "+rule.Name)
	}
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatList(ruleList))
	_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")
}
