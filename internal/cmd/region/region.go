// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package region

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "region",
		Usage:           "Administer, detail, and interact with Attila regions",
		HideHelpCommand: true,
		UsageText:       "attila region <command> [options] [args]",
		Subcommands: []*cli.Command{
			createCommand(),
			deleteCommand(),
			getCommand(),
			listCommand(),
			shellCommand(),
		},
	}
}

func outputRegion(cliCtx *cli.Context, r *api.Region) {

	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
		fmt.Sprintf("Name|%s", r.Name),
		fmt.Sprintf("Group|%s", r.Group),
		fmt.Sprintf("TLS Enabled|%v", r.TLS != nil),
		fmt.Sprintf("Create Time|%s", helper.FormatTime(r.Metadata.CreateTime)),
		fmt.Sprintf("Update Time|%s", helper.FormatTime(r.Metadata.UpdateTime)),
	}))

	out := make([]string, 0, len(r.API)+1)
	out = append(out, "Address|Default")

	for _, apiEndpoint := range r.API {
		out = append(out, fmt.Sprintf(
			"%s|%v", apiEndpoint.Address, apiEndpoint.Default,
		))
	}

	_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n\n")
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatList(out))
	_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")

	if r.TLS != nil {
		fmt.Println("\nTLS CA Cert")
		fmt.Println(r.TLS.CACert)
		fmt.Println("\nTLS Client Cert")
		fmt.Println(r.TLS.ClientCert)
		fmt.Println("\nTLS Client Key")
		fmt.Println(r.TLS.ClientKey)
	}
}
