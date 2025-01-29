SOURCE_IMAGE_REPO ?= gcr.io/k8s-staging-sig-storage/objectstorage-sidecar
SOURCE_IMAGE_VERSION ?= v20250123-a0e4046
SOURCE_IMAGE ?= $(SOURCE_IMAGE_REPO):$(SOURCE_IMAGE_VERSION)

TARGET_IMAGE_REPO ?= ghcr.io/mesosphere/dkp-container-images/objectstorage-sidecar
TARGET_IMAGE_VERSION ?= $(SOURCE_IMAGE_VERSION)
TARGET_IMAGE ?= $(TARGET_IMAGE_REPO):$(TARGET_IMAGE_VERSION)

.PHONY: docker-build
docker-build:
	docker build --build-arg="SOURCE_IMAGE=$(SOURCE_IMAGE)" -t $(TARGET_IMAGE) .

.PHONY: build-args
build-args:
	@echo "SOURCE_IMAGE=$(SOURCE_IMAGE)"
	@echo "TARGET_IMAGE=$(TARGET_IMAGE)"
	@echo "TARGET_IMAGE_VERSION=$(TARGET_IMAGE_VERSION)"
