# CNPG PostgreSQL Container Images

A custom build of `ghcr.io/cloudnative-pg/postgresql` container image.
Includes `pgvector` extension, which is not included by default in the upstream `minimal` image.
The `minimal` flavour is used to reduce the amount of CVEs.

## Build

```shell
make docker-build
```
