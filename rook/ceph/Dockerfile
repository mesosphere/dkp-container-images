# syntax=docker/dockerfile:1

ARG SOURCE_IMAGE
FROM $SOURCE_IMAGE

USER root
RUN yum remove python3-joblib -y
USER 2016
