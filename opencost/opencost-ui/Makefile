# OpenCost UI - Build from Source
# Production builds should use the build-from-source.yaml GitHub workflow

SOURCE_REPO ?= opencost/opencost-ui
SOURCE_VERSION ?= 1.117.5
UI_PATH ?= /dkp/opencost/frontend

TARGET_IMAGE_REPO ?= ghcr.io/mesosphere/dkp-container-images/opencost/opencost-ui
TARGET_IMAGE ?= $(TARGET_IMAGE_REPO):$(SOURCE_VERSION)

.PHONY: build-args
build-args:
	@echo "ui_path=$(UI_PATH)"
	@echo "SOURCE_REPO=$(SOURCE_REPO)"
	@echo "SOURCE_VERSION=$(SOURCE_VERSION)"
	@echo "TARGET_IMAGE=$(TARGET_IMAGE)"
