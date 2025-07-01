# syntax=docker/dockerfile:1
ARG SOURCE_IMAGE

FROM ${SOURCE_IMAGE}

ARG PG_MAJOR_VERSION

USER root

RUN apt-get update && \
	apt-get install -y --no-install-recommends postgresql-${PG_MAJOR_VERSION}-pgvector && \
    apt-get clean && \
	rm -rf /var/lib/apt/lists/*

USER postgres
