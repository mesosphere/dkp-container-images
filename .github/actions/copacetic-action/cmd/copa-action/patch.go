package copa_action

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/goware/prefixer"
	"github.com/spf13/cobra"

	"github.com/d2iq-labs/copacetic-action/pkg/cli"
	"github.com/d2iq-labs/copacetic-action/pkg/image"
	"github.com/d2iq-labs/copacetic-action/pkg/patch"
	"github.com/d2iq-labs/copacetic-action/pkg/registry"
)

var (
	ghcrOrg = "mesosphere/dkp-container-images"

	// <original>-d2iq.<version>
	imageTagSuffix = "d2iq"
	debug          = false
	timeout        = 1 * time.Hour
	skipUpload     = false
	format         = "json"
)

func NewPatchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "patch PATH | -",
		Short: "Patch runs copatetic patch operation on list of provided images",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithTimeout(cmd.Context(), timeout)
			defer cancel()

			input, err := cli.OpenFileOrStdin(args[0])
			if err != nil {
				return err
			}

			images, err := cli.ReadImages(input)
			if err != nil {
				return err
			}

			ghcr := registry.NewGHCR(ghcrOrg)
			if skipUpload {
				ghcr.WithSkipUploads(slog.Default())
			}

			tasks := []*patch.Task{}
			for _, imageRef := range images {
				slog.Info("patching image", "imageRef", imageRef)
				logger, _ := getImageLogger(imageRef, debug)
				task, err := patch.Run(ctx, imageRef, ghcr, imageTagSuffix, debug, logger)
				if err != nil {
					cmdErr := &image.CmdErr{}
					if errors.As(err, &cmdErr) {
						slog.Error("failed patching", "err", err)
						io.Copy(os.Stderr, prefixer.New(bytes.NewReader(cmdErr.Output), "output | "))
					} else {
						slog.Error("failed patching", "err", err)
					}
				}
				tasks = append(tasks, task)
			}

			return patch.WriteJSON(tasks, os.Stdout)
		},
	}

	cmd.Flags().BoolVar(&debug, "debug", debug, "enable debugging")
	cmd.Flags().BoolVar(&skipUpload, "skip-upload", skipUpload, "skip uploading to remote registry")
	cmd.Flags().DurationVar(&timeout, "timeout", timeout, "total timeout of run")
	cmd.Flags().StringVar(&ghcrOrg, "ghcr-org", ghcrOrg, "name of ghcr.io org where patched images will get published")
	cmd.Flags().StringVar(&imageTagSuffix, "patched-tag-suffix", imageTagSuffix, "name of patched images tag suffix, e.g v1.0.0 will become v1.0.0-d2iq.0")

	return cmd
}

func getImageLogger(imageRef string, debug bool) (*slog.Logger, io.Reader) {
	data := &bytes.Buffer{}
	var imagePatchLogs io.Writer
	imagePatchLogs = data
	if debug {
		imagePatchLogs = io.MultiWriter(data, os.Stderr)
	}

	logger := slog.New(slog.NewTextHandler(
		imagePatchLogs, &slog.HandlerOptions{})).With("imageRef", imageRef)
	return logger, data
}
