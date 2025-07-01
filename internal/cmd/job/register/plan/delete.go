// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package plan

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/pkg/api"
)

func deleteCommand() *cli.Command {
	return &cli.Command{
		Name:      "delete",
		Usage:     "Delete a job registration plan",
		Category:  "plan",
		Args:      true,
		UsageText: "attila job register plan delete [options] [plan-id]",
		Flags:     helper.ClientFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					deleteCLIErrorMsg,
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			id, err := ulid.Parse(cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError(deleteCLIErrorMsg, err), 1)
			}

			client := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			deleteReq := api.JobRegisterPlanDeleteReq{ID: id}

			_, err = client.JobRegisterPlans().Delete(context.Background(), &deleteReq)
			if err != nil {
				return cli.Exit(helper.FormatError(deleteCLIErrorMsg, err), 1)
			}

			_, _ = fmt.Fprintf(cliCtx.App.Writer, "successfully deleted job registration plan %q\n", id)
			return nil
		},
	}
}
