package registry

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/github"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/daemon"
	"github.com/google/go-containerregistry/pkg/v1/remote/transport"
)

var ErrImageNotFound = errors.New("image not found")

// Registry where the patched images are stored.
type Registry interface {
	// ListTags returns list of patched tags for given source image.
	ListTags(ctx context.Context, sourceImageRef string) ([]string, error)

	// ImageRef returns name of given image in the patched registry.
	ImageRef(sourceImage, tag string) (string, error)

	// OriginalImageRef returns original image for given patched image.
	OriginalImageRef(imageRef string) string

	// Pushes container image from local daemon to registry.
	Push(ctx context.Context, daemonImage, targetImage string) error
}

const ghcrDomain = "ghcr.io"

// NewGHCR creates ghcr.io registry for storing patched images under a single
// organization in ghcr.
func NewGHCR(organization string) *ghcrRegistry {
	return &ghcrRegistry{
		organization: organization,
		keychain:     github.Keychain,
		skipUpload:   false,
	}
}

type ghcrRegistry struct {
	keychain     authn.Keychain
	organization string
	skipUpload   bool
	logger       *slog.Logger
}

var _ Registry = &ghcrRegistry{}

func (r *ghcrRegistry) WithSkipUploads(logger *slog.Logger) {
	r.skipUpload = true
	r.logger = logger
}

func (r *ghcrRegistry) ListTags(ctx context.Context, sourceImageRef string) ([]string, error) {
	ref, err := name.ParseReference(sourceImageRef)
	if err != nil {
		return nil, err
	}

	registryImageName := r.patchedImageName(ref)
	tags, err := crane.ListTags(
		registryImageName,
		r.defaultCraneOpts(ctx)...,
	)
	if r.isNotFoundError(err) {
		return nil, ErrImageNotFound
	}

	return tags, err
}

func (r *ghcrRegistry) patchedImageName(ref name.Reference) string {
	imageRef := ref.Context().String()
	return fmt.Sprintf("%s/%s/%s", ghcrDomain, r.organization, r.fixIndexDockerRegistry(imageRef))
}

func (r *ghcrRegistry) fixIndexDockerRegistry(imageRef string) string {
	if strings.HasPrefix(imageRef, "index.docker.io") {
		return strings.TrimPrefix(imageRef, "index.")
	}
	return imageRef
}

func (r *ghcrRegistry) defaultCraneOpts(ctx context.Context) []crane.Option {
	return []crane.Option{
		crane.WithContext(ctx),
		crane.WithAuthFromKeychain(r.keychain),
	}
}

func (r *ghcrRegistry) isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	transportError := &transport.Error{}
	if errors.As(err, &transportError) {
		if transportError.StatusCode == http.StatusForbidden {
			return true
		}
	}

	if strings.Contains(err.Error(), "NAME_UNKNOWN: repository name not known to registry") {
		return true
	}

	return false
}

// ImageRef returns reference of patched image in the cache registry.
func (r *ghcrRegistry) ImageRef(sourceImageRef, tag string) (string, error) {
	ref, err := name.ParseReference(sourceImageRef, name.WithDefaultRegistry("docker.io"))
	if err != nil {
		return "", err
	}

	taggedRef := r.fixIndexDockerRegistry(ref.Context().Tag(tag).String())
	return fmt.Sprintf("%s/%s/%s", ghcrDomain, r.organization, taggedRef), nil
}

// OriginalImageRef returns name of imageRef from which was the given imageRef built.
func (r *ghcrRegistry) OriginalImageRef(imageRef string) string {
	registryPrefix := fmt.Sprintf("%s/%s/", ghcrDomain, r.organization)
	if !strings.HasPrefix(imageRef, registryPrefix) {
		return ""
	}

	imageRefWithoutPrefix := strings.TrimPrefix(imageRef, registryPrefix)
	if !strings.Contains(imageRefWithoutPrefix, "/") {
		return ""
	}

	baseRef, _ := name.ParseReference(imageRefWithoutPrefix, name.WithDefaultRegistry("docker.io"))
	baseTag := ParseBaseTag(baseRef.Identifier())

	return baseRef.Context().Tag(baseTag).String()
}

func (r *ghcrRegistry) Push(ctx context.Context, daemonImage, targetImage string) error {
	if r.skipUpload {
		r.logger.Info("skipping image push", "daemonImage", daemonImage, "targetImage", targetImage)
		return nil
	}

	daemonImageRef, err := name.ParseReference(daemonImage)
	if err != nil {
		return err
	}

	patchedImage, err := daemon.Image(daemonImageRef, daemon.WithContext(ctx))
	if err != nil {
		return err
	}

	// Upload image to registry
	if err := crane.Push(patchedImage, targetImage, r.defaultCraneOpts(ctx)...); err != nil {
		return err
	}

	return nil
}
