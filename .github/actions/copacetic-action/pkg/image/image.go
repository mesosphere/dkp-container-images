package image

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"

	"github.com/d2iq-labs/copacetic-action/pkg/registry"
)

// assumption
//   - always rebuild image from source image
// options
// 1. original image with no patched version in registry (foo/nxingx:v1)
//    - scan original and check if there are packages for update
//    - build new patched image
//    - publish to patch registry
//    - report: upgrade possible from original => patched
// 2. original image with patched version in registry (ghcr.io/mesosphere/foo/nginx:v1-d2iq.XXX)
//    - scan patched and check if there are packages for update
//    - if no new fixed packages
//       - report: upgrade possible from original to patched
//    - else build new patched image and publish
//       - report upgrade possible from original to new patched image
// 3. latest patched image in registry
// 4. patched image is used (v1-d2iq.(XXX-n)in registry has newever image (v1-d2iq.(XXX-n))

// input imageRef
//   - is original
//   - get original (for patching only) - input for patching
//   - get latest patch - input for scanning

// scanned image (original or latest patch)
//   if there are fixable vulnerabilites patch original
//   publish patched image
//   store patched image

// report phase
// input image => there is patched with higher version

// ImagePatch represents a state for patching a single image.
type ImagePatch struct {
	// By default same as input or source image from which the Input was built from by
	// previous patching.
	Source string

	sourceRef name.Reference

	// Image that is used for scanning. Latest patched image or original.
	Scanned string

	// Name of image that should be used instead of Input image if patching was
	// successful.
	// By default its the latest patched version in the registry. If patching process
	// generated a new image this value will be modified.
	// Used for reporting purposes to output which images should be replaced.
	Patched string

	tagResolver  registry.TagResolver
	existingTags []string
}

func (i *ImagePatch) NextPatchedTag() string {
	return i.tagResolver.Next(i.existingTags)
}

// SourceRef is a parsed reference of source image that was used for building
// the patched image.
func (i *ImagePatch) SourceRef() name.Reference {
	return i.sourceRef
}

func NewImagePatch(
	ctx context.Context,
	imageRef string,
	reg registry.Registry,
	imageTagSuffix string,
) (*ImagePatch, error) {
	patch := &ImagePatch{
		Source:  imageRef,
		Scanned: imageRef,
	}

	// Calculate source reference if provided image is already patched by previous runs.
	if originalImageRef := reg.OriginalImageRef(imageRef); originalImageRef != "" {
		patch.Source = originalImageRef
	}

	var err error
	patch.sourceRef, err = name.ParseReference(patch.Source)
	if err != nil {
		return nil, fmt.Errorf("failed to parse source image ref %q: %w", patch.Source, err)
	}

	// List existing tags for the original image.
	patch.existingTags, err = reg.ListTags(ctx, patch.Source)
	if err != nil && !errors.Is(err, registry.ErrImageNotFound) {
		return nil, fmt.Errorf("failed to list tags for %q: %w", patch.Source, err)
	}

	patch.tagResolver = registry.NewTagResolver(patch.sourceRef.Identifier(), imageTagSuffix)

	// Check find which image tag is latest in patched repository, which will be scanned.
	if latesImageTag := patch.tagResolver.Latest(patch.existingTags); latesImageTag != "" {
		latestImageRef, err := reg.ImageRef(patch.Source, latesImageTag)
		if err != nil {
			return nil, fmt.Errorf("failed to generate latest image ref for %q: %w", patch.Source, err)
		}
		patch.Scanned = latestImageRef
	}

	// Scanned is the latest build patched image. If there already exists image that is patched
	// use it as default result for given scanned image.
	if imageRef != patch.Scanned {
		patch.Patched = patch.Scanned
	}

	return patch, nil
}
