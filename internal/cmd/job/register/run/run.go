// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

const (
	createCLIErrorMsg = "failed to create job registration run"
	deleteCLIErrorMsg = "failed to delete job registration run"
	getCLIErrorMsg    = "failed to get job registration run"
	listCLIErrorMsg   = "failed to list job registration run"
)

func Command() *cli.Command {
	return &cli.Command{
		Name:            "run",
		Category:        "register",
		Usage:           "View and manage Attila job registration runs",
		HideHelpCommand: true,
		UsageText:       "attila job register run <command> [options] [args]",
		Subcommands: []*cli.Command{
			createCommand(),
			deleteCommand(),
			getCommand(),
			listCommand(),
		},
	}
}

func outputRuns(cliCtx *cli.Context, runs []*api.JobRegisterRunStub) {
	if len(runs) == 0 {
		_, _ = fmt.Fprint(cliCtx.App.Writer, "No Attila job registration runs found")
		return
	}

	out := make([]string, 0, len(runs)+1)
	out = append(out, "ID|Namespace|Job ID|Job Namespace")
	for _, run := range runs {
		out = append(out, fmt.Sprintf(
			"%s|%s|%s|%s",
			run.ID, run.Namespace, run.JobID, run.JobNamespace))
	}

	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatList(out))
	_, _ = fmt.Fprint(cliCtx.App.Writer, "\n")
}

func outputRun(cliCtx *cli.Context, run *api.JobRegisterRun) {
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
		fmt.Sprintf("ID|%s", run.ID),
		fmt.Sprintf("Namespace|%s", run.Namespace),
		fmt.Sprintf("Num Regions|%v", len(run.Regions)),
		fmt.Sprintf("Job ID|%s", run.JobID),
		fmt.Sprintf("Job Namespace|%s", run.JobNamespace),
	}))

	for _, regionRun := range run.Regions {
		var errStr string
		if regionRun.Error != nil {
			errStr = regionRun.Error.Error()
		}

		var evalID string
		if regionRun.RegisterResp != nil {
			evalID = regionRun.RegisterResp.EvalID
		}

		_, _ = fmt.Fprint(cliCtx.App.Writer, color.New(color.Bold).Sprintf(
			"\n\nRegion %q Run:\n", regionRun.Region))

		_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
			fmt.Sprintf("Eval ID|%s", evalID),
			fmt.Sprintf("Warnings|%s", errorString(regionRun.RegisterResp.Warnings)),
			fmt.Sprintf("Error|%s", errStr),
		}))

		_, _ = fmt.Fprint(cliCtx.App.Writer, "\n")
	}
}

func errorString(e any) string {
	switch v := e.(type) {
	case error:
		if v != nil {
			return v.Error()
		}
	case string:
		return v
	default:
		return ""
	}
	return ""
}
