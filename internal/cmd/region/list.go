// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package region

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
		Usage:     "List Attila regions",
		Category:  "region",
		Args:      false,
		UsageText: "attila region list [options]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			regions, _, err := client.Regions().List(context.Background())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to list Attila regions", err), 1)
			}

			_, _ = fmt.Fprint(cliCtx.App.Writer, formatRegionList(regions.Regions))
			_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")

			return nil
		},
	}
}

func formatRegionList(regions []*api.RegionStub) string {
	if len(regions) == 0 {
		return "No Attila Regions found"
	}

	out := make([]string, 0, len(regions)+1)
	out = append(out, "Name|Group|TLS|Addresses")
	for _, region := range regions {
		out = append(out, fmt.Sprintf(
			"%s|%s|%v|%v",
			region.Name, region.Group, region.TLSEnabled, strings.Join(region.Addresses, ", ")))
	}

	return helper.FormatList(out)
}
