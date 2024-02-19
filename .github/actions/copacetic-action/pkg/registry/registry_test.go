package registry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegistryOriginalImageRef(t *testing.T) {
	registry := NewGHCR("d2iq-labs")

	testCases := []struct {
		imageRef string
		expected string
	}{
		{"mesosphere/kommander-cm:v2.0.0", ""},
		{"ghcr.io/mesosphere/kommander-cm:v2.0.0", ""},
		{"ghcr.io/d2iq-labs/kommander-cm:v2.0.0", ""},
		{"ghcr.io/d2iq-labs/mesosphere/kommander-cm:v2.0.0-d2iq.2", "mesosphere/kommander-cm:v2.0.0"},
		{
			"ghcr.io/d2iq-labs/registry.k8s.io/sig-storage/local-volume-provisioner:v2.5.0-d2iq.0",
			"registry.k8s.io/sig-storage/local-volume-provisioner:v2.5.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.imageRef, func(t *testing.T) {
			assert.Equal(t, tc.expected, registry.OriginalImageRef(tc.imageRef))
		})
	}
}

func TestRegistryListTags(t *testing.T) {
	r := NewGHCR("d2iq-labs")
	tags, err := r.ListTags(context.Background(), "registry.k8s.io/sig-storage/local-volume-provisioner:v2.5.0")
	assert.ErrorIs(t, err, ErrImageNotFound)
	assert.Empty(t, tags)
}
