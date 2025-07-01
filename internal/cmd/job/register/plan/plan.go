// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package plan

import (
	"fmt"
	"sort"
	"strings"

	"github.com/fatih/color"
	nomadAPI "github.com/hashicorp/nomad/api"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

const (
	createCLIErrorMsg = "failed to create job registration plan"
	deleteCLIErrorMsg = "failed to delete job registration plan"
	getCLIErrorMsg    = "failed to get job registration plan"
	listCLIErrorMsg   = "failed to list job registration plans"
	runCLIErrorMsg    = "failed to run job registration plan"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "plan",
		Category:        "register",
		Usage:           "Plan and execute Nomad job registrations",
		HideHelpCommand: true,
		UsageText:       "attila job register plan <command> [options] [args]",
		Subcommands: []*cli.Command{
			createCommand(),
			deleteCommand(),
			getCommand(),
			listCommand(),
			runCommand(),
		},
	}
}

func outputPlan(cliCtx *cli.Context, plan *api.JobRegisterPlan) {
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
		fmt.Sprintf("ID|%s", plan.ID),
		fmt.Sprintf("Num Regions|%v", len(plan.Regions)),
		fmt.Sprintf("Job ID|%s", plan.JobID),
		fmt.Sprintf("Job Namespace|%s", plan.JobNamespace),
	}))
	_, _ = fmt.Fprint(cliCtx.App.Writer, "\n")

	for _, regionPlan := range plan.Regions {
		outputPlannedJob(cliCtx, regionPlan)
	}
}

func outputPlannedJob(cliCtx *cli.Context, regionPlan *api.JobRegisterRegionPlan) {

	type outDetail struct {
		sortedTG []string
		kv       map[string][]string
		list     map[string][]string
	}

	out := outDetail{
		kv:       make(map[string][]string),
		list:     make(map[string][]string),
		sortedTG: sortedTaskGroupNames(regionPlan.Plan.Annotations.DesiredTGUpdates),
	}

	for _, tg := range out.sortedTG {

		allocUpdates := regionPlan.Plan.Annotations.DesiredTGUpdates[tg]

		out.kv[tg] = append(out.kv[tg], []string{
			fmt.Sprintf("Ignored Allocations|%v", allocUpdates.Ignore),
			fmt.Sprintf("Placed Allocations|%v", allocUpdates.Place),
			fmt.Sprintf("Migrated Allocations|%v", allocUpdates.Migrate),
			fmt.Sprintf("Stopped Allocations|%v", allocUpdates.Stop),
			fmt.Sprintf("In-place Updated Allocations|%v", allocUpdates.InPlaceUpdate),
			fmt.Sprintf("Destroyed Allocations|%v", allocUpdates.DestructiveUpdate),
			fmt.Sprintf("Canary Allocations|%v", allocUpdates.Canary),
			fmt.Sprintf("Preempted Allocations|%v", allocUpdates.Preemptions),
		}...)
	}

	for _, tg := range sortedTaskGroupNames(regionPlan.Plan.FailedTGAllocs) {
		kv, outList := failedAllocDetail(regionPlan.Plan.FailedTGAllocs[tg])

		out.kv[tg] = append(out.kv[tg], kv...)
		out.list[tg] = append(out.list[tg], outList...)
	}

	for _, tg := range out.sortedTG {

		_, _ = fmt.Fprint(cliCtx.App.Writer, color.New(color.Bold).Sprintf(
			"\nRegion %q Plan for Task Group %q:\n", regionPlan.Region, tg))

		if kv, ok := out.kv[tg]; ok {
			_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV(kv))
			_, _ = fmt.Fprint(cliCtx.App.Writer, "\n\n")
		}

		if list, ok := out.list[tg]; ok && len(list) > 1 {
			_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatList(list))
			_, _ = fmt.Fprint(cliCtx.App.Writer, "\n")
		}
	}
}

func failedAllocDetail(metrics *nomadAPI.AllocationMetric) ([]string, []string) {
	outKV := []string{
		fmt.Sprintf("Allocation Placement Failures|%v", metrics.CoalescedFailures+1),
		fmt.Sprintf("Nodes Evaluated|%v", metrics.NodesEvaluated),
		fmt.Sprintf("Nodes Exhausted|%v", metrics.NodesExhausted),
	}

	for dc, numNodes := range metrics.NodesAvailable {
		outKV = append(outKV, fmt.Sprintf("Nodes Available In Datacenter %q|%v", dc, numNodes))
	}

	outKV = append(outKV, fmt.Sprintf("Quotas Exhauted|%s", strings.Join(metrics.QuotaExhausted, ", ")))

	//
	outList := []string{"Type|Detail|Num Nodes"}

	for class, num := range metrics.ClassFiltered {
		outList = append(outList, fmt.Sprintf("%s|%s|%v", "class filtered", class, num))
	}
	for class, num := range metrics.ConstraintFiltered {
		outList = append(outList, fmt.Sprintf("%s|%s|%v", "constraint filtered", class, num))
	}
	for class, num := range metrics.ClassExhausted {
		outList = append(outList, fmt.Sprintf("%s|%s|%v", "class exhausted", class, num))
	}
	for dim, num := range metrics.DimensionExhausted {
		outList = append(outList, fmt.Sprintf("%s|%s|%v", "dimension exhausted", dim, num))
	}

	return outKV, outList
}

func sortedTaskGroupNames[V any](groups map[string]V) []string {
	tgs := make([]string, 0, len(groups))
	for tg := range groups {
		tgs = append(tgs, tg)
	}
	sort.Strings(tgs)
	return tgs
}
