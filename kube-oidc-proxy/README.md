# kube-oidc-proxy

The original `kube-oidc-proxy` project haven't been updated in years. That means that libraries and official container image now includes vulnerable code.

There is a [maintained fork](https://www.tremolosecurity.com/post/updating-kube-oidc-proxy) which publishes updated version. The container image uses `ubuntu` as a base image, which means it cannot be included in the DKP airgapped bundle, due to Ubuntu licensing.

The forked image gets rebuilt by copying the fork build and adding it to static distroless container image to minimize attack surface.

## Build

```
make docker-build SOURCE_IMAGE_VERSION=1.x.x
```

