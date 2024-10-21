BASE_IMAGE ?= gcr.io/distroless/static-debian12@sha256:69830f29ed7545c762777507426a412f97dad3d8d32bae3e74ad3fb6160917ea

SOURCE_IMAGE_REPO ?= ghcr.io/tremolosecurity/kube-oidc-proxy
SOURCE_IMAGE_VERSION ?= 1.0.7
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
