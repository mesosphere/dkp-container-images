# OpenCost UI - Build from Source
# Production builds should use the build-from-source.yaml GitHub workflow

SOURCE_REPO ?= opencost/opencost-ui
SOURCE_VERSION ?= ${SOURCE_VERSION?"Need to set SOURCE_VERSION"}
UI_PATH ?= /dkp/opencost/frontend
SOURCE_REPO_PATH ?= ../../source-repo

TARGET_IMAGE_REPO ?= ghcr.io/mesosphere/dkp-container-images/opencost/opencost-ui
TARGET_IMAGE ?= $(TARGET_IMAGE_REPO):$(SOURCE_VERSION:v%=%)

.PHONY: build-args
build-args: patch-upstream
	@echo "ui_path=$(UI_PATH)"

.PHONY: image-values
image-values:
	@echo "TARGET_IMAGE=$(TARGET_IMAGE)"

.PHONY: patch-upstream:
patch-upstream:
	@FILE="$(SOURCE_REPO_PATH)/default.nginx.conf.template"; \
	if [ -f "$$FILE" ]; then \
		echo "Found Nginx template. Removing IPv6 listener..."; \
		sed -i '/listen \[::\]/d' "$$FILE"; \
		grep "listen" "$$FILE" || true; \
	else \
		echo "File $$FILE not found. Skipping patch."; \
	fi
