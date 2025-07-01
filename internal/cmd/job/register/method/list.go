// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package method

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
		Usage:     "List Attila job registration methods",
		Category:  "method",
		Args:      false,
		UsageText: "attila job register method list [options]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			methodListResp, _, err := client.JobRegisterMethods().List(context.Background())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to list Attila job registration methods", err), 1)
			}

			_, _ = fmt.Fprint(cliCtx.App.Writer, formatMethodList(methodListResp.Methods))
			_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")

			return nil
		},
	}
}

func formatMethodList(methods []*api.JobRegisterMethodStub) string {
	if len(methods) == 0 {
		return "No Attila job registration methods found"
	}

	out := make([]string, 0, len(methods)+1)
	out = append(out, "Name|Selector")
	for _, method := range methods {
		out = append(out, fmt.Sprintf(
			"%s|%s",
			method.Name, method.Selector))
	}

	return helper.FormatList(out)
}
