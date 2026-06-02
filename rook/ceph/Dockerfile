# syntax=docker/dockerfile:1

ARG SOURCE_IMAGE
ARG JQ_IMAGE

FROM ${JQ_IMAGE} as jq
FROM ${SOURCE_IMAGE}

USER root
RUN dnf remove python3-joblib jq -y
COPY --from=jq --chmod=755 /jq /usr/bin/jq
USER 2016
