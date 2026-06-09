SOURCE_IMAGE_REPO ?= quay.io/ceph/ceph
SOURCE_IMAGE_VERSION ?= v19.2.4
SOURCE_IMAGE ?= $(SOURCE_IMAGE_REPO):$(SOURCE_IMAGE_VERSION)

JQ_IMAGE_REPO ?= ghcr.io/jqlang/jq
JQ_IMAGE_VERSION ?= 1.8.2rc1
JQ_IMAGE ?= $(JQ_IMAGE_REPO):$(JQ_IMAGE_VERSION)

REPO ?= $(shell gh repo view --json name --template '{{.name}}')
ORG ?= $(shell gh repo view --json owner --template '{{.owner.login}}')

TARGET_IMAGE_REPO ?= ghcr.io/${ORG}/${REPO}/ceph/ceph
TARGET_IMAGE_VERSION ?= $(SOURCE_IMAGE_VERSION)
TARGET_IMAGE ?= $(TARGET_IMAGE_REPO):$(TARGET_IMAGE_VERSION)

.PHONY: docker-build
docker-build:
	docker build --build-arg="SOURCE_IMAGE=$(SOURCE_IMAGE)" --build-arg="JQ_IMAGE=$(JQ_IMAGE)" -t $(TARGET_IMAGE) .

.PHONY: build-args
build-args:
	@echo "SOURCE_IMAGE=$(SOURCE_IMAGE)"
	@echo "JQ_IMAGE=$(JQ_IMAGE)"
	@echo "TARGET_IMAGE=$(TARGET_IMAGE)"
	@echo "TARGET_IMAGE_VERSION=$(TARGET_IMAGE_VERSION)"
