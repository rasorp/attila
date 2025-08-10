// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package topology

import (
	"fmt"

	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func listCommand() *cli.Command {
	return &cli.Command{
		Name:      "list",
		Usage:     "List an overview of the collected topologies",
		Category:  "topology",
		Args:      true,
		UsageText: "attila topology list [options]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			topologies, _, err := client.Topologies().List(cliCtx.Context, nil)
			if err != nil {
				return cli.Exit(helper.FormatError("failed to list Attila topologies", err), 1)
			}

			_, _ = fmt.Fprint(cliCtx.App.Writer, formatTopologyList(topologies.TopologyOverviews))
			_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")

			return nil
		},
	}
}

func formatTopologyList(overviews []*api.TopologyOverview) string {
	if len(overviews) == 0 {
		return "No Topologies found"
	}

	out := make([]string, 0, len(overviews)+1)
	out = append(
		out,
		"Region Name|Num Servers|Num Clients|Num Allocs|CPU (MHz)|Memory (MB)",
	)

	for _, overview := range overviews {
		out = append(out, fmt.Sprintf(
			"%s|%v|%v|%v|%v|%v",
			overview.RegionName,
			overview.NumServers,
			overview.NumClients,
			overview.NumAllocs,
			fmt.Sprintf("%v/%v", overview.CPUAllocated, overview.CPUAllocatable),
			fmt.Sprintf("%v/%v", overview.MemoryAllocated, overview.MemoryAllocatable),
		))
	}

	return helper.FormatList(out)
}
