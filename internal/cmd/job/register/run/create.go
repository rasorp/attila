// Copyright James Rasell 2025, 2026
// SPDX-License-Identifier: Apache-2.0

package run

import (
	"fmt"
	"os"

	"github.com/hashicorp/nomad/jobspec2"
	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func createCommand() *cli.Command {
	return &cli.Command{
		Name:      "create",
		Usage:     "Create an Attila job registration run",
		Category:  "run",
		Args:      true,
		UsageText: "attila job register run create [options]",
		Flags:     append(createFlags(), helper.ClientNamespaceFlags()...),
		Action: func(cliCtx *cli.Context) error {

			planID := cliCtx.String("plan-id")
			jobspecFile := ""

			// Positional argument is the jobspec file.
			if n := cliCtx.Args().Len(); n == 1 {
				jobspecFile = cliCtx.Args().First()
			} else if n > 1 {
				return cli.Exit(helper.FormatError(createCLIErrorMsg,
					fmt.Errorf("expected at most 1 argument (jobspec file), got %v", n)), 1)
			}

			// Exactly one of plan-id or jobspec-file must be provided.
			hasPlanID := planID != ""
			hasJobspecFile := jobspecFile != ""
			if !hasPlanID && !hasJobspecFile {
				return cli.Exit(helper.FormatError(createCLIErrorMsg,
					fmt.Errorf("either --plan-id or a jobspec file must be provided")), 1)
			}
			if hasPlanID && hasJobspecFile {
				return cli.Exit(helper.FormatError(createCLIErrorMsg,
					fmt.Errorf("only one of --plan-id or jobspec file can be provided")), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			var createReq *api.JobsRegisterRunsCreateReq

			if hasPlanID {
				ulid, err := ulid.Parse(planID)
				if err != nil {
					return cli.Exit(helper.FormatError(createCLIErrorMsg,
						fmt.Errorf("failed to parse plan-id: %w", err)), 1)
				}
				createReq = &api.JobsRegisterRunsCreateReq{PlanID: ulid}
			} else {
				// Parse the jobspec file.
				jobspecBytes, err := os.ReadFile(jobspecFile)
				if err != nil {
					return cli.Exit(helper.FormatError(createCLIErrorMsg,
						fmt.Errorf("failed to read jobspec file: %w", err)), 1)
				}

				parsedJobspec, err := jobspec2.ParseWithConfig(&jobspec2.ParseConfig{
					Path:     jobspecFile,
					Body:     jobspecBytes,
					ArgVars:  cliCtx.StringSlice("jobspec-var"),
					VarFiles: cliCtx.StringSlice("jobspec-var-file"),
					Strict:   true,
				})
				if err != nil {
					return cli.Exit(helper.FormatError(createCLIErrorMsg,
						fmt.Errorf("failed to parse jobspec: %w", err)), 1)
				}

				createReq = &api.JobsRegisterRunsCreateReq{Job: parsedJobspec}
			}

			resp, _, err := client.JobRegisterRuns().Create(
				cliCtx.Context,
				createReq,
				&api.WriteOpts{Namespace: helper.NamespaceFromFlags(cliCtx)},
			)
			if err != nil {
				return cli.Exit(helper.FormatError(createCLIErrorMsg, err), 1)
			}

			outputRun(cliCtx, resp.Run)
			return nil
		},
	}
}

func createFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  "plan-id",
			Value: "",
			Usage: "The plan ID to use for the registration run",
		},
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
