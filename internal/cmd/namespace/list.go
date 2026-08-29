// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package namespace

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List Attila namespaces",
		Category:  "namespace",
		Args:      false,
		UsageText: "attila namespace list [options]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			nsResp, _, err := client.Namespaces().List(cliCtx.Context)
			if err != nil {
				return cli.Exit(helper.FormatError("failed to list Attila namespaces", err), 1)
			}

			_, _ = fmt.Fprint(cliCtx.App.Writer, formatNamespaceList(nsResp.Namespaces))
			_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")

			return nil
		},
	}
}

func formatNamespaceList(namespaces []*api.NamespaceStub) string {
	if len(namespaces) == 0 {
		return "No Attila namespaces found"
	}

	out := make([]string, 0, len(namespaces)+1)
	out = append(out, "Name|Description")
	for _, ns := range namespaces {
		out = append(out, fmt.Sprintf("%s|%s", ns.Name, ns.Description))
	}

	return helper.FormatList(out)
}
