# dkp-container-images

Custom and maintained container images used by **NKP** (Nutanix Kubernetes Platform). This repo was originally created for DKP (D2iQ Kubernetes Platform); DKP is now **NKP**. The repo builds, patches, and publishes images for security, licensing, and feature customizations (e.g. airgapped bundles, distroless bases, extra extensions).

## Is this repo active?

**Yes.** The repository is actively maintained. As a guide, consider it active if there has been a check-in within the last few months. At the time of this note, the last commit was **2025-12-30** (merge of the create-publish-doca-ofed workflow). Check `git log -1` or the GitHub commit history for the current last commit.

## Where images are pushed

All images are published to **GitHub Container Registry (GHCR)** under the org/repo path:

- **Registry:** `ghcr.io`
- **Image path pattern:** `ghcr.io/mesosphere/dkp-container-images/<image-name>`

Examples:

| Image | Full image reference |
|-------|------------------------|
| Ceph | `ghcr.io/mesosphere/dkp-container-images/ceph/ceph:<version>` |
| Rook Ceph | `ghcr.io/mesosphere/dkp-container-images/rook/ceph:<version>` |
| CloudNative-PG PostgreSQL | `ghcr.io/mesosphere/dkp-container-images/cloudnative-pg/postgresql:<version>` |
| OpenCost UI | `ghcr.io/mesosphere/dkp-container-images/opencost/opencost-ui:<version>` |
| kube-oidc-proxy | `ghcr.io/mesosphere/dkp-container-images/kube-oidc-proxy:<version>` |
| COSI controller/sidecar | `ghcr.io/mesosphere/dkp-container-images/objectstorage-controller:<version>` etc. |

CI logs in to GHCR and pushes only when the workflow input **Push to registry** is set to `true` (default is `false` for safety).

## Usage in NKP: management and workload clusters

**Yes.** Both the **management** and **workload** clusters use images with the prefix `ghcr.io/mesosphere/dkp-container-images`.

**Management cluster** (e.g. `dm-nkp-mgmt-1.conf`) and **workload cluster** (e.g. `dm-nkp-workload-1.kubeconfig`) use the same set of mirrored images. Example images in use:

| Image | Full image reference |
|-------|------------------------|
| CloudNative-PG PostgreSQL | `ghcr.io/mesosphere/dkp-container-images/cloudnative-pg/postgresql:17.5-minimal-bookworm` |
| kube-oidc-proxy | `ghcr.io/mesosphere/dkp-container-images/kube-oidc-proxy:1.0.9` |
| COSI objectstorage-controller | `ghcr.io/mesosphere/dkp-container-images/objectstorage-controller:v20250110-a29e5f6` |
| COSI objectstorage-sidecar | `ghcr.io/mesosphere/dkp-container-images/objectstorage-sidecar:v20240513-v0.1.0-35-gefb3255` |
| OpenCost UI | `ghcr.io/mesosphere/dkp-container-images/opencost/opencost-ui:1.118.0` |

The same registry prefix also appears in:

- **Docs:** logging (Fluent Bit), Prometheus/Alertmanager (Karma), and related config examples.
- **Support bundle** (`troubleshoot.sh/support-bundle-...`): Fluent Bit, Rook/Ceph, Karma, Kubecost frontend, kube-oidc-proxy, PostgreSQL, Flux kustomize-controller, Weaviate, nginx-unprivileged, and others.

So **DKP’s mirrored images from `ghcr.io/mesosphere/dkp-container-images` are in use on both management and workload clusters** (and in docs and support-bundle references).

---

## Adding a new image

Choose one of two patterns depending on whether the image is built from an **in-repo Dockerfile** or from **upstream source**.

### Option A: In-repo Dockerfile (rebuild / repackage upstream)

Use this when you have a **Dockerfile in this repo** that uses an upstream image as base or copies artifacts (e.g. Ceph, kube-oidc-proxy, COSI, CloudNative-PG).

1. **Create a directory** (e.g. `my-component/my-image` or `my-image` at top level).

2. **Add at minimum:**
   - `Dockerfile` – build that uses build args for source image/version.
   - `Makefile` with:
     - Variables: `SOURCE_IMAGE` (or equivalent), `TARGET_IMAGE`, `TARGET_IMAGE_VERSION`.
     - `docker-build` target that builds using those variables.
     - **`build-args` target** that prints one `KEY=VALUE` per line for every build-arg and `TARGET_IMAGE` (and `TARGET_IMAGE_VERSION` if used). The **Rebuild image** workflow relies on this.

3. **Set the target registry in the Makefile:**

   ```makefile
   TARGET_IMAGE_REPO ?= ghcr.io/mesosphere/dkp-container-images/<path>/<name>
   TARGET_IMAGE ?= $(TARGET_IMAGE_REPO):$(TARGET_IMAGE_VERSION)
   ```

4. **Add a short README** in that directory describing what the image is and how to build it (see existing `*/README.md` under `ceph/`, `kube-oidc-proxy/`, `cosi/`, `cloudnative-pg/`).

5. **Build and push via GitHub Actions:**
   - **Actions → Rebuild image**
   - **directory:** path to the image directory (e.g. `kube-oidc-proxy` or `ceph/ceph`).
   - **build-args:** optional, e.g. `SOURCE_IMAGE_VERSION=1.0.9`.
   - **platforms:** default `linux/amd64`; add more for multi-arch if needed.
   - **push:** set to `true` when ready to publish to GHCR.

6. **Local test (optional):**  
   From the image directory: `make docker-build` (override vars as needed). Push manually to GHCR if you’re not using the workflow yet.

### Option B: Build from upstream source

Use this when the image is built by **cloning an upstream repo** and building there (e.g. OpenCost UI).

1. **Create a directory** under this repo (e.g. `opencost/opencost-ui`) that will **only** hold the Makefile and README (no Dockerfile; the Dockerfile lives in the upstream repo).

2. **Add a Makefile** that provides:
   - **`build-args`** – prints build args for the upstream Dockerfile (e.g. `ui_path=...`), one `KEY=VALUE` per line.
   - **`image-values`** – prints `TARGET_IMAGE=<full tag>` (and any other vars the workflow expects). The workflow runs `make build-args` and `make image-values` from `./automation-repo/<source-repo>` (e.g. `opencost/opencost-ui`), so paths must match the **source-repo** input.

3. **Set the target image** to `ghcr.io/mesosphere/dkp-container-images/<path>/<name>:<version>` in the Makefile.

4. **Add a README** explaining the image, the upstream repo, and that production builds use **Build image from source**.

5. **Build and push via GitHub Actions:**
   - **Actions → Build image from source**
   - **source-repo:** upstream repo (e.g. `opencost/opencost-ui`) — must match the directory path under this repo used in the workflow.
   - **source-version:** tag/branch/commit (e.g. `v1.117.5`).
   - **build-args:** optional; extra args in `KEY=VALUE` form.
   - **platforms:** default `linux/amd64`.
   - **push:** set to `true` to push to GHCR.

---

## Updating an existing image

- **In-repo Dockerfile:**  
  - Bump the version (e.g. `SOURCE_IMAGE_VERSION`) in the image’s **Makefile** (or pass it via workflow **build-args**).  
  - Run **Actions → Rebuild image** with the correct **directory** and **push: true** when ready.

- **Build from source:**  
  - Run **Actions → Build image from source** with a new **source-version** (and any **build-args**). Set **push: true** to publish.  
  - Optionally update the README or default version in the Makefile in this repo.

In both cases, the image is pushed to **`ghcr.io/mesosphere/dkp-container-images/...`** when **push** is enabled.

---

## CVE patching

The **CVE patch** workflow (Actions → **CVE patch**) rebuilds and patches existing images for CVEs and pushes the patched images. It can be triggered manually (input: list of **images**) or via `repository_dispatch` with type `patch-images`. Patched images are published to GHCR (under the same `ghcr.io/mesosphere/...` naming used by this repo).

---

## Summary

| Goal | Workflow | Where it pushes |
|------|----------|-----------------|
| New/updated image from in-repo Dockerfile | **Rebuild image** | `ghcr.io/mesosphere/dkp-container-images/<directory>` |
| New/updated image from upstream source | **Build image from source** | `ghcr.io/mesosphere/dkp-container-images/<source-repo path>` |
| Patch existing images for CVEs | **CVE patch** | Same GHCR locations |

Per-image details (build args, versions, extensions) are documented in each component’s README (e.g. `ceph/ceph/README.md`, `kube-oidc-proxy/README.md`, `opencost/opencost-ui/README.md`).
