package image

import (
	"context"
	"fmt"
	"log/slog"
	"path"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
)

func AddLabels(ctx context.Context, patchedRef name.Reference, labels map[string]string) error {
	patchedImage, err := daemon.Image(patchedRef, daemon.WithContext(ctx))
	if err != nil {
		return err
	}

	patchedImageConfig, err := patchedImage.ConfigFile()
	if err != nil {
		return err
	}
	patchedImageConfig = patchedImageConfig.DeepCopy()

	if patchedImageConfig.Config.Labels == nil {
		patchedImageConfig.Config.Labels = map[string]string{}
	}
	for k, v := range labels {
		patchedImageConfig.Config.Labels[k] = v
	}
	patchedImage, err = mutate.Config(patchedImage, patchedImageConfig.Config)
	if err != nil {
		return err
	}

	_, err = daemon.Write(patchedRef.(name.Tag), patchedImage)
	return err
}

func PatchCVEs(ctx context.Context, imageRef, buildTag, tmpDir string) error {
	report, err := Scan(ctx, imageRef)
	if err != nil {
		return err
	}

	reportPath := path.Join(tmpDir, "copa-report.json")
	if err := report.WriteTo(reportPath); err != nil {
		return fmt.Errorf("failed to write trivy report for %q: %w", imageRef, err)
	}

	output, err := copa(
		ctx,
		"-i", imageRef,
		"patch", "-r", reportPath,
		"-t", buildTag,
		"--timeout", "5m",
		// copa locks if error is produced during the package validation
		// See: https://github.com/project-copacetic/copacetic/issues/503
		"--ignore-errors",
		"--debug",
	)
	if err != nil {
		return &CmdErr{
			Err:    err,
			Output: output,
		}
	}

	return err
}

func copa(ctx context.Context, args ...string) ([]byte, error) {
	cmd, _, _ := prepareCmd(ctx, "copa", args...)
	cmd.Stdout = nil
	cmd.Stderr = nil

	output, err := cmd.CombinedOutput()
	if err != nil {
		slog.ErrorContext(ctx, "failed to exec cmd", "cmdName", "copa", "err", err)
	}
	slog.Debug("command completed", "cmdName", "copa", "args", args, "output", output)
	return output, err
}
