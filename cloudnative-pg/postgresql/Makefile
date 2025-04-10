SOURCE_IMAGE_REPO ?= ghcr.io/cloudnative-pg/postgresql
SOURCE_IMAGE_VERSION ?= 17.4-minimal-bookworm
SOURCE_IMAGE ?= $(SOURCE_IMAGE_REPO):$(SOURCE_IMAGE_VERSION)
PG_MAJOR_VERSION ?= $(firstword $(subst ., ,$(SOURCE_IMAGE_VERSION)))

TARGET_IMAGE_REPO ?= ghcr.io/mesosphere/dkp-container-images/cloudnative-pg/postgresql
TARGET_IMAGE_VERSION ?= $(SOURCE_IMAGE_VERSION)
TARGET_IMAGE ?= $(TARGET_IMAGE_REPO):$(TARGET_IMAGE_VERSION)

.PHONY: docker-build
docker-build:
	docker build -t $(TARGET_IMAGE) \
		--build-arg="SOURCE_IMAGE=$(SOURCE_IMAGE)" \
		--build-arg="PG_MAJOR_VERSION=$(PG_MAJOR_VERSION)" \
		$(CURDIR)

.PHONY: build-args
build-args:
	@echo "SOURCE_IMAGE=$(SOURCE_IMAGE)"
	@echo "PG_MAJOR_VERSION=$(PG_MAJOR_VERSION)"
	@echo "TARGET_IMAGE=$(TARGET_IMAGE)"
	@echo "TARGET_IMAGE_VERSION=$(TARGET_IMAGE_VERSION)"
