// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"fmt"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "namespace",
		Usage:           "Administer, detail, and interact with Attila namespaces",
		HideHelpCommand: true,
		UsageText:       "attila namespace <command> [options] [args]",
		Subcommands: []*cli.Command{
			createCommand(),
			deleteCommand(),
			getCommand(),
			listCommand(),
		},
	}
}

func outputNamespace(cliCtx *cli.Context, ns *api.Namespace) {
	outputKV := []string{
		fmt.Sprintf("Name|%s", ns.Name),
		fmt.Sprintf("Description|%s", ns.Description),
		fmt.Sprintf("Regions|%s", strings.Join(ns.Regions, ", ")),
	}

	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV(outputKV))
}
