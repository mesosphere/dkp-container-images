BASE_IMAGE ?= gcr.io/distroless/static-debian12@sha256:0d6ada5faececed5cd3f99baa08e4109934f2371c0d81b3bff38924fe1deea05

SOURCE_IMAGE_REPO ?= ghcr.io/tremolosecurity/kube-oidc-proxy
SOURCE_IMAGE_VERSION ?= 1.0.6
SOURCE_IMAGE ?= $(SOURCE_IMAGE_REPO):$(SOURCE_IMAGE_VERSION)

TARGET_IMAGE_REPO ?= ghcr.io/mesosphere/dkp-container-images/kube-oidc-proxy
TARGET_IMAGE_VERSION ?= $(SOURCE_IMAGE_VERSION)
TARGET_IMAGE ?= $(TARGET_IMAGE_REPO):$(TARGET_IMAGE_VERSION)

.PHONY: docker-build
docker-build:
	docker build --build-arg="BASE_IMAGE=$(BASE_IMAGE)" --build-arg="SOURCE_IMAGE=$(SOURCE_IMAGE)" -t $(TARGET_IMAGE) .

.PHONY: build-args
build-args:
	@echo "BASE_IMAGE=$(BASE_IMAGE)"
	@echo "SOURCE_IMAGE=$(SOURCE_IMAGE)"
	@echo "TARGET_IMAGE=$(TARGET_IMAGE)"
	@echo "TARGET_IMAGE_VERSION=$(TARGET_IMAGE_VERSION)"
