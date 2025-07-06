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

	outputKV := []string{
		fmt.Sprintf("Name|%s", r.Name),
		fmt.Sprintf("Group|%s", r.Group),
		fmt.Sprintf("TLS Enabled|%v", r.TLS != nil),
	}

	if r.TLS != nil {
		outputKV = append(outputKV,
			fmt.Sprintf("TLS Server Name|%s", r.TLS.ServerName),
			fmt.Sprintf("TLS Insecure|%v", r.TLS.Insecure),
		)
	}

	outputKV = append(outputKV,
		fmt.Sprintf("Create Time|%s", helper.FormatTime(r.Metadata.CreateTime)),
		fmt.Sprintf("Update Time|%s", helper.FormatTime(r.Metadata.UpdateTime)),
	)

	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV(outputKV))

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
		_, _ = fmt.Fprintf(cliCtx.App.Writer, "\nTLS CA Cert:")
		_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n%s", r.TLS.CACert)
		_, _ = fmt.Fprintf(cliCtx.App.Writer, "\nTLS Client Cert:")
		_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n%s", r.TLS.ClientCert)
		_, _ = fmt.Fprintf(cliCtx.App.Writer, "\nTLS Client Key:")
		_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n%s", r.TLS.ClientKey)
	}
}
