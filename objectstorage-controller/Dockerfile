# syntax=docker/dockerfile:1

# From image should be `gcr.io/k8s-staging-sig-storage/objectstorage-controller`.
ARG SOURCE_IMAGE

FROM ${SOURCE_IMAGE}

# Same ENTRYPOINT as in the SOURCE_IMAGE
ENTRYPOINT ["/controller"]
