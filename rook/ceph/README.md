# rook/ceph

A custom build of `rook/ceph` container image that removes the `python-joblib` which is installed as part of Python scikit and is only required for disk usage prediction, which is not enabled by default.

## Build

```shell
make docker-build
```
