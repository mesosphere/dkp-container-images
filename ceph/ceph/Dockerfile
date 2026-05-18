# syntax=docker/dockerfile:1

ARG SOURCE_IMAGE
ARG JQ_IMAGE

FROM ${JQ_IMAGE} as jq
FROM ${SOURCE_IMAGE}

RUN dnf remove python3-joblib jq -y

COPY --from=jq /jq /usr/bin/jq
RUN chmod +x /usr/bin/jq
USER ceph
