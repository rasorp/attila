// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package topology

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func getCommand() *cli.Command {
	return &cli.Command{
		Name:      "get",
		Usage:     "Get detailed topology for a region",
		Category:  "topology",
		Args:      true,
		UsageText: "attila topology get [options] [region-name]",
		Flags:     getFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					"failed to get Attila topology",
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			topology, _, err := client.Topologies().Get(
				cliCtx.Context,
				&api.TopologiesGetReq{RegionName: cliCtx.Args().First()},
			)
			if err != nil {
				return cli.Exit(helper.FormatError("failed to get Attila topology", err), 1)
			}

			outputTopology(cliCtx, topology.Topology)
			return nil
		},
	}
}

// getFlags returns the full set of CLI flags for use with the topology get
// command including the API client set.
func getFlags() []cli.Flag {
	return append(helper.ClientFlags(), &cli.BoolFlag{
		Name:  "node-allocs",
		Usage: "Include node allocation topology in the output",
		Value: false,
	})
}

func outputTopology(cliCtx *cli.Context, topology *api.Topology) {

	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
		fmt.Sprintf("Region Name|%s", topology.Overview.RegionName),
		fmt.Sprintf("Num Servers|%v", topology.Overview.NumServers),
		fmt.Sprintf("Num Clients|%v", topology.Overview.NumClients),
		fmt.Sprintf("Num Allocs|%v", topology.Overview.NumAllocs),
		fmt.Sprintf("CPU MHz|%v/%v", topology.Overview.CPUAllocated, topology.Overview.CPUAllocatable),
		fmt.Sprintf("Memory MB|%v/%v", topology.Overview.MemoryAllocated, topology.Overview.MemoryAllocatable),
		fmt.Sprintf("Create Time|%v", topology.CreateTime),
	}))
	_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n\n")

	serverList := make([]string, 0, len(topology.Detail.Servers)+1)
	serverList = append(serverList, "ID|Name|Status|Version|Raft Version")
	for _, server := range topology.Detail.Servers {
		serverList = append(serverList, fmt.Sprintf(
			"%s|%s|%s|%s|%s",
			server.ID, server.Name, server.Status, server.Version, server.RaftVersion))
	}

	_, _ = fmt.Fprint(cliCtx.App.Writer, color.New(color.Bold).Sprintf("Server Topology:\n"))
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatList(serverList))

	// Sort the nodes listing by node pool then node name, before outputting the
	// detail to the console. This provides a consistent view and makes it
	// easier to grok on each invocation.
	topology.Detail.SortNodes()
	outputNodeTopology(cliCtx, topology.Detail.Nodes)
}

func outputNodeTopology(cliCtx *cli.Context, nodeTopology []*api.NodeTopology) {

	// Identify if the caller wants to include the allocation topology in the
	// output.
	includeAllocs := cliCtx.Bool("node-allocs")

	// Build the list header which has conditional headings depending on if the
	// node allocations should be included.
	header := "ID|Name|Node Pool|Status"

	if includeAllocs {
		header += "|Alloc ID|Job ID|Namespace"
	}

	header += "|CPU MHz|Memory MB"

	nodeListOutput := []string{header}

	for _, node := range nodeTopology {

		nodeOutput := fmt.Sprintf(
			"%s|%s|%s|%s",
			node.ID,
			node.Name,
			node.NodePool,
			node.Status,
		)

		// If we are including allocation detail, mark the allocation specific
		// headings with asterisks for the node row.
		if includeAllocs {
			nodeOutput += fmt.Sprintf("|%s|%s|%s", "*", "*", "*")
		}

		// Add the node CPU and memory heading which details the allocatable vs.
		// allocated resources.
		nodeOutput += fmt.Sprintf(
			"|%v/%v|%v/%v",
			node.CPUAllocated,
			node.CPUAllocatable,
			node.MemoryAllocated,
			node.MemoryAllocatable,
		)

		nodeListOutput = append(nodeListOutput, nodeOutput)

		if includeAllocs {

			// Sort the node allocations, so we have a consistent view on each
			// run of the CLI.
			node.SortAllocs()

			// Iterate the node allocations and add a row for each one.
			// Asterisks are used for columns that represent information that
			// only applies to the node object.
			for _, alloc := range node.AllocationTopology {
				nodeListOutput = append(
					nodeListOutput,
					fmt.Sprintf(
						"%s|%s|%s|%s|%s|%s|%s|%v|%v",
						"*",
						"*",
						"*",
						"*",
						alloc.ID,
						alloc.JobID,
						alloc.Namespace,
						alloc.CPU,
						alloc.Memory,
					),
				)
			}
		}
	}

	_, _ = fmt.Fprint(cliCtx.App.Writer, color.New(color.Bold).Sprintf("\n\nNode Topology:\n"))
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatList(nodeListOutput))
	_, _ = fmt.Fprint(cliCtx.App.Writer, "\n")
}
