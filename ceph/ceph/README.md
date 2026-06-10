# quay.io/ceph/ceph

A custom build of `quay.io/ceph/ceph` container image that removes the `python-joblib` which is installed as part of Python scikit and is only required for disk usage prediction, which is not enabled by default and upgrades the jq utility to version `1.8.2rc1`

## Build

```shell
make docker-build
```
