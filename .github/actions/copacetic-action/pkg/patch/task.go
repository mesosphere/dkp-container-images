package patch

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"slices"

	"github.com/google/go-containerregistry/pkg/name"
	"go.step.sm/crypto/randutil"

	"github.com/d2iq-labs/copacetic-action/pkg/image"
	"github.com/d2iq-labs/copacetic-action/pkg/registry"
)

type Task struct {
	Error error
	Patch *image.ImagePatch
	Image string
}

func Run(ctx context.Context, imageRef string, reg registry.Registry, imageTagSuffix string, debug bool, logger *slog.Logger) (*Task, error) {
	parsed, err := name.ParseReference(imageRef)
	if err != nil {
		return nil, err
	}
	if _, ok := parsed.(name.Digest); ok {
		return nil, fmt.Errorf("images with digest references are not supported")
	}

	t := &Task{
		Image: imageRef,
	}

	// Create new image patch record.
	imagePatch, err := image.NewImagePatch(ctx, imageRef, reg, imageTagSuffix)
	if err != nil {
		return nil, fmt.Errorf("failed to create image patch for %q: %w", imageRef, err)
	}
	logger.Info("processing image", "image", imagePatch)

	t.Patch = imagePatch

	withErr := func(t *Task, err error) *Task {
		t.Error = err
		return t
	}

	// To avoid generating same patched image always scan the latest patched
	// tag in the patch registry and only build image if there are available
	// fixes that would change the latest patched version.
	report, err := image.Scan(ctx, imagePatch.Scanned)
	if err != nil {
		return withErr(t, err), err
	}

	if len(report.Vulnerabilities()) == 0 {
		logger.Info("no fixable vulnerabilities found in scanned image", "scannedImage", imagePatch.Scanned)
		return t, nil
	}

	logger.Info(
		"found patchable vulnerabilities",
		"scanned", imagePatch.Scanned,
		"vulnerabilites", report.Vulnerabilities(),
	)

	buildId, err := randutil.Alphanumeric(5)
	logger.Info("generated unique buildId", "buildId", buildId)
	if err != nil {
		return withErr(t, err), err
	}

	tmpDir, err := os.MkdirTemp("", fmt.Sprintf("patch-%s-", buildId))
	if err != nil {
		return withErr(t, err), err
	}
	logger.Info("created tmp dir", "tmpDir", tmpDir)
	defer func() {
		if debug {
			logger.Warn("debug enabled workdir not cleaned up", "path", tmpDir)
			return
		}
		os.RemoveAll(tmpDir)
	}()

	// Patch image
	buildTag := fmt.Sprintf("copa-patched-%s", buildId)
	err = image.PatchCVEs(ctx, imagePatch.Source, buildTag, tmpDir)
	if err != nil {
		return withErr(t, err), err
	}

	// Local image name built by copacetic
	patchedRef := imagePatch.SourceRef().Context().Tag(buildTag)
	logger.Info("regenerated image using copa", "patchedRef", patchedRef.String())

	patchedReport, err := image.Scan(ctx, patchedRef.String())
	if err != nil {
		return withErr(t, err), err
	}
	logger.Info(
		"patched vulnerability report",
		"original", report.Vulnerabilities(),
		"patched", patchedReport.Vulnerabilities(),
	)

	if slices.Equal(
		image.VulnerabilitiesIdsSorted(report.Vulnerabilities()),
		image.VulnerabilitiesIdsSorted(patchedReport.Vulnerabilities()),
	) {
		logger.Warn("no vulnerabilties were fixed by running copa",
			"scannedImage", imagePatch.Scanned,
			"scanned", image.VulnerabilitiesIdsSorted(report.Vulnerabilities()),
			"patched", image.VulnerabilitiesIdsSorted(patchedReport.Vulnerabilities()),
		)
		return t, nil
	}

	// Add labels to the newly built image
	labels := map[string]string{
		"com.d2iq.source-image": imagePatch.Source,
	}
	if err := image.AddLabels(ctx, patchedRef, labels); err != nil {
		return withErr(t, err), err
	}

	// Push patched image from local docker daemon to the remote registry.
	{
		imagePatch.Patched, err = reg.ImageRef(imagePatch.Source, imagePatch.NextPatchedTag())
		if err != nil {
			return withErr(t, err), err
		}

		logger.Info("uploading image to registry", "patchedImage", imagePatch.Patched)
		if err := reg.Push(ctx, patchedRef.String(), imagePatch.Patched); err != nil {
			return withErr(t, err), err
		}
	}

	return t, nil
}
