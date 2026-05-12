# syntax=docker/dockerfile:1

ARG SOURCE_IMAGE
ARG JQ_IMAGE

FROM ${JQ_IMAGE} as jq
FROM $SOURCE_IMAGE

RUN yum remove python3-joblib -y

USER root
COPY --from=jq /jq /usr/bin/jq
RUN chmod +x /usr/bin/jq
USER ceph
