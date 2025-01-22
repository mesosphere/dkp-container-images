# objectstorage-controller

A custom build of `gcr.io/k8s-staging-sig-storage/objectstorage-controller:` container image. The registry hosting this image is being shut down https://console.cloud.google.com/gcr/images/k8s-staging-sig-storage/global/objectstorage-controller. 

The Dockerile will be based on the upstream project https://github.com/kubernetes-sigs/container-object-storage-interface/blob/main/controller/Dockerfile.

## Build

```shell
make docker-build
```
