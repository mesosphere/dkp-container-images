# OpenCost UI - Build from Source
# Production builds should use the build-from-source.yaml GitHub workflow

SOURCE_REPO ?= opencost/opencost-ui
SOURCE_VERSION ?= ${SOURCE_VERSION?"Need to set SOURCE_VERSION"}
UI_PATH ?= /dkp/opencost/frontend

REPO ?= $(shell gh repo view --json name --template '{{.name}}')
ORG ?= $(shell gh repo view --json owner --template '{{.owner.login}}')

TARGET_IMAGE_REPO ?= ghcr.io/${ORG}/${REPO}/opencost/opencost-ui
TARGET_IMAGE ?= $(TARGET_IMAGE_REPO):$(SOURCE_VERSION:v%=%)

.PHONY: build-args
build-args:
	@echo "ui_path=$(UI_PATH)"

.PHONY: image-values
image-values:
	@echo "TARGET_IMAGE=$(TARGET_IMAGE)"
