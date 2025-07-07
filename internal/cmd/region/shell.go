// Copyright (c) James Rasell
// SPDX-License-Identifier: Apache-2.0

package region

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/urfave/cli/v2"

	"github.com/rasorp/attila/internal/cmd/helper"
	"github.com/rasorp/attila/internal/helper/pointer"
	"github.com/rasorp/attila/pkg/api"
)

func shellCommand() *cli.Command {
	return &cli.Command{
		Name:      "shell",
		Usage:     "Run and enter a regional shell",
		Category:  "region",
		Args:      true,
		UsageText: "attila region shell [options] [region-name]",
		Flags:     shellFlags(),
		Action: func(cliCtx *cli.Context) error {

			if numArgs := cliCtx.Args().Len(); numArgs != 1 {
				return cli.Exit(helper.FormatError(
					"failed to run region shell",
					fmt.Errorf("expected 1 argument, got %v", numArgs)),
					1,
				)
			}

			attilaClient := api.NewClient(helper.ClientConfigFromFlags(cliCtx))

			regionResp, _, err := attilaClient.Regions().Get(context.Background(), cliCtx.Args().First())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to run region shell", err), 1)
			}

			dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
			if err != nil {
				return cli.Exit(helper.FormatError("failed to run region shell", err), 1)
			}

			defer func(dockerClient *client.Client) { _ = dockerClient.Close() }(dockerClient)

			ctx := context.Background()

			// Build our Docker image reference using the CLI flag variable as
			// the version identifier. In the future, we may want to allow full
			// configuration of the image path, so we account for private
			// registry use.
			dockerRef := "hashicorp/nomad:" + cliCtx.String("nomad-image-version")

			reader, err := dockerClient.ImagePull(ctx, dockerRef, image.PullOptions{})
			if err != nil {
				return cli.Exit(helper.FormatError("failed to run region shell", err), 1)
			}

			// If the returned pull reader is not nil, then the pull is in
			// progress. Wait until we receive to EOF to indicate the pull is
			// complete, otherwise the pull has failed.
			if reader != nil {
				defer func(reader io.ReadCloser) { _ = reader.Close() }(reader)
				_, err = io.Copy(io.Discard, reader)
				if err != nil && !errors.Is(err, io.EOF) {
					return cli.Exit(helper.FormatError("failed to run region shell", err), 1)
				}
			}

			resp, err := dockerClient.ContainerCreate(ctx, &container.Config{
				Image:        dockerRef,
				Entrypoint:   []string{""},
				Cmd:          []string{"/bin/ash"},
				AttachStdout: true,
				OpenStdin:    true,
				Tty:          true,
				Hostname:     "attila-region-shell-" + regionResp.Region.Name,
				Env:          buildEnvVarList(regionResp.Region),
			}, nil, nil, nil, "")
			if err != nil {
				return cli.Exit(helper.FormatError("failed to run region shell", err), 1)
			}

			if err := dockerClient.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
				return cli.Exit(helper.FormatError("failed to run region shell", err), 1)
			}

			attachResp, err := dockerClient.ContainerAttach(ctx, resp.ID, container.AttachOptions{
				Stdout: true,
				Stdin:  true,
				Stream: true,
			})
			if err != nil {
				panic(err)
			}

			// Start Go routines to copy the containers output to the terminal,
			// and copy the terminal input into the container.
			go func() { _, _ = io.Copy(os.Stdout, attachResp.Reader) }()
			go func() { _, _ = io.Copy(attachResp.Conn, os.Stdin) }()

			// Write a newline to the container terminal, so the local terminal
			// shows us the correct command prompt.
			_, _ = attachResp.Conn.Write([]byte("\n"))

			// Wait for the container to reach a non-running state, or for the
			// operator to interrupt the program and want to quit.
			statusCh, errCh := dockerClient.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)

			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

			select {
			case <-sigs:
			case <-statusCh:
			case <-errCh:
			}

			attachResp.Close()

			if err := dockerClient.ContainerStop(ctx, resp.ID, container.StopOptions{Timeout: pointer.Of(0)}); err != nil {
				return cli.Exit(helper.FormatError("failed to run region shell", err), 1)
			}
			if err := dockerClient.ContainerRemove(ctx, resp.ID, container.RemoveOptions{Force: true}); err != nil {
				return cli.Exit(helper.FormatError("failed to run region shell", err), 1)
			}

			return nil
		},
	}
}

// shellFlags returns the full set of CLI flags for use with the region shell
// command including the API client set.
func shellFlags() []cli.Flag {
	return append(helper.ClientFlags(), &cli.StringFlag{
		Name:  "nomad-image-version",
		Usage: "The Nomad Docker image version identifier to use",
		Value: "1.10.2",
	})
}

func buildEnvVarList(r *api.Region) []string {
	out := []string{
		"NOMAD_ADDR=" + r.DefaultOrFirstAddress(),
		"NOMAD_REGION=" + r.Name,
	}

	if r.TLS != nil {
		out = append(out,
			"NOMAD_CA_CERT="+r.TLS.CACert,
			"NOMAD_CLIENT_CERT="+r.TLS.ClientCert,
			"NOMAD_CLIENT_KEY="+r.TLS.ClientKey,
			"NOMAD_TLS_SERVER_NAME="+r.TLS.ServerName,
			"NOMAD_SKIP_VERIFY="+strconv.FormatBool(r.TLS.Insecure),
		)
	}

	return out
}
