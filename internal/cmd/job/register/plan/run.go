// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package plan

import (
	"context"
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hashicorp/nomad/jobspec2"
	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func runCommand() *cli.Command {
	return &cli.Command{
		Name:      "run",
		Usage:     "Run a job registration plan",
		Category:  "plan",
		Args:      true,
		UsageText: "attila job register plan run [options] [plan-id] [job-spec]",
		Flags:     append(helper.ClientFlags(), runFlags()...),
		Action: func(cliCtx *cli.Context) error {

			cliArgs := cliCtx.Args()

			if cliArgs.Len() != 2 {
				return cli.Exit(helper.FormatError(runCLIErrorMsg,
					fmt.Errorf("expected 2 arguments, got %v", cliArgs.Len())), 1)
			}

			id, err := ulid.Parse(cliArgs.First())
			if err != nil {
				return cli.Exit(helper.FormatError(runCLIErrorMsg, err), 1)
			}

			jobfile := cliArgs.Get(1)

			jobsepcBytes, err := os.ReadFile(jobfile)
			if err != nil {
				return cli.Exit(helper.FormatError(runCLIErrorMsg,
					fmt.Errorf("failed to read jobspec file: %w", err)), 1)
			}

			jobspecParseConfig := jobspec2.ParseConfig{
				Path:     jobfile,
				Body:     jobsepcBytes,
				ArgVars:  cliCtx.StringSlice("jobspec-var"),
				VarFiles: cliCtx.StringSlice("jobspec-var-file"),
				Strict:   true,
			}

			parsedJobspec, err := jobspec2.ParseWithConfig(&jobspecParseConfig)
			if err != nil {
				return cli.Exit(helper.FormatError(runCLIErrorMsg,
					fmt.Errorf("failed to parse jobspec: %w", err)), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			req := api.JobsRegisterPlanRunReq{ID: id, Job: parsedJobspec}

			resp, _, err := client.JobRegisterPlans().Run(context.Background(), &req)
			if err != nil {
				return cli.Exit(helper.FormatError(runCLIErrorMsg, err), 1)
			}

			outputPlanRun(cliCtx, resp)
			return nil
		},
	}
}

func runFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringSliceFlag{
			Name:  "jobspec-var-file",
			Value: cli.NewStringSlice(),
			Usage: "The path to a HCL2 file containing user variables",
		},
		&cli.StringSliceFlag{
			Name:  "jobspec-var",
			Value: cli.NewStringSlice(),
			Usage: "A HCL2 user variable",
		},
	}
}

func outputPlanRun(cliCtx *cli.Context, resp *api.JobsRegisterPlanRunResp) {
	_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
		fmt.Sprintf("ID|%s", resp.Run.ID),
		fmt.Sprintf("Num Regions|%v", len(resp.Run.Regions)),
		fmt.Sprintf("Job ID|%s", resp.Run.JobID),
		fmt.Sprintf("Job Namespace|%s", resp.Run.JobNamespace),
		fmt.Sprintf("Partial Error|%s", errorString(resp.PatrialFailureError)),
	}))

	for _, regionPlan := range resp.Run.Regions {
		_, _ = fmt.Fprint(cliCtx.App.Writer, color.New(color.Bold).Sprintf(
			"\n\nRegion %q Run:\n", regionPlan.Region))

		_, _ = fmt.Fprint(cliCtx.App.Writer, helper.FormatKV([]string{
			fmt.Sprintf("Eval ID|%s", regionPlan.Run.EvalID),
			fmt.Sprintf("Warnings|%s", regionPlan.Run.Warnings),
			fmt.Sprintf("Error|%s", errorString(regionPlan.Error)),
		}))

		_, _ = fmt.Fprintf(cliCtx.App.Writer, "\n")
	}
}

func errorString(e error) string {
	if e != nil {
		return e.Error()
	}
	return ""
}
