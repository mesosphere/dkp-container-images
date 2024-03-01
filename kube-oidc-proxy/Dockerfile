# syntax=docker/dockerfile:1

# From image should be `kube-oidc-proxy` binary copied to the new base image.
ARG SOURCE_IMAGE
# Name of the base image to which should be the `kube-oidc-proxy` binary copied.
ARG BASE_IMAGE

FROM $SOURCE_IMAGE as source

FROM ${BASE_IMAGE}

# RUN apk --no-cache add ca-certificates && update-ca-certificates

# See: https://github.com/TremoloSecurity/kube-oidc-proxy/blob/1acdbd93710e61117d6daf90f541f2107552d9e6/Dockerfile
COPY --from=source /usr/bin/kube-oidc-proxy /usr/local/bin/kube-oidc-proxy

CMD [ "/usr/local/bin/kube-oidc-proxy" ]
