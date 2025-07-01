// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package plan

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/nomad/jobspec2"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create a job registration plan",
		Category:  "plan",
		Args:      true,
		UsageText: "attila job register plan create [options] [job-spec]",
		Flags:     append(helper.ClientFlags(), createFlags()...),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					createCLIErrorMsg,
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			jobfile := cliCtx.Args().First()

			jobsepcBytes, err := os.ReadFile(jobfile)
			if err != nil {
				return cli.Exit(helper.FormatError(
					createCLIErrorMsg,
					fmt.Errorf("failed to read jobspec file: %w", err)),
					1,
				)
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
				return cli.Exit(helper.FormatError(
					createCLIErrorMsg,
					fmt.Errorf("failed to parse jobspec: %w", err)),
					1,
				)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			req := api.JobRegisterPlanCreateReq{Job: parsedJobspec}

			resp, _, err := client.JobRegisterPlans().Create(context.Background(), &req)
			if err != nil {
				return cli.Exit(helper.FormatError(createCLIErrorMsg, err), 1)
			}

			outputPlan(cliCtx, resp.Plan)
			return nil
		},
	}
}

func createFlags() []cli.Flag {
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
