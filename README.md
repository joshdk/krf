[![License][license-badge]][license-link]
[![Actions][github-actions-badge]][github-actions-link]
[![Releases][github-release-badge]][github-release-link]

# Kubernetes Resource Filter

🔬 Command line utility for filtering collections of Kubernetes resources

## What is this?

This repository provides `krf`, a tool that can be used to easily sift through large collections of Kubernetes resources using a number of powerful filters.

The original impetus for this was the need to search through a stream of Kubernetes yaml documents for resources that contain some value. 
The immediate next problem is that tools like e.g. `grep` are only filename/line number aware at best; so you end up with a list of matching lines, but you are no closer to knowing which specific resources contain those lines.
To address this, `krf` was build as a Kubernetes resource aware `grep`, but can be utilized with both more power and precision.

## Installation

### Release artifact

Binaries for various architectures are published on the [releases][github-release-link] page.

The latest release can be installed by running:

```shell
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m | sed 's/x86_64/amd64/; s/aarch64/arm64/')
curl -Lso krf.tar.gz https://github.com/joshdk/krf/releases/latest/download/krf-${OS}-${ARCH}.tar.gz
tar -xf krf.tar.gz
sudo mkdir -p /usr/local/bin/
sudo install krf /usr/local/bin/
```

### Brew

Release binaries are also available via [Brew](https://brew.sh).

The latest release can be installed by running:

```shell
brew tap joshdk/tap
brew install joshdk/tap/krf
```

### Go install

Installation can also be done directly from this repository.

The latest commit can be installed by running:

```shell
go install github.com/joshdk/krf@master
```

## Usage

### Reading Resources

To begin, a corpus of input resources can be provided in a number of ways.

A single file:
```shell
krf resources.yaml
```

All yaml files contained within a directory:
```shell
krf ./manifests/
```

The output of `kustomize` build:
```shell
kustomize build … | krf
```

The output of `kubectl` in yaml or json format:
```shell
kubectl get … -o=yaml | krf
kubectl get … -o=json | krf
```

### Filtering Resources

The input corpus can then be filtered using a set of individual _matchers_ that you can mix and match.

Include resources that have the name `backend`:
```shell
… | krf --name backend
```

Include resources that have the name `backend` and the namespace `default`:
```shell
… | krf --name backend --namespace default
```

Include resources that have the name `backend` and are of kind `Deployment` or `Service`:
```shell
… | krf --name backend --kind deploy,svc
```

Include resources that have the name `backend` or `frontend` and are not of kind `Ingress`:
```shell
… | krf --name backend,frontend --not-kind ing
```

Finally, simply output the same set of resources that were provided as input:
```shell
… | krf
```

#### Filtering Logic

When `krf` is invoked, each input resource is evaluated against various categories of positive and negative matchers in order to reject, or ultimately accept the resource.

At a high level, in order to be accepted, each resource:
- **Must not** match **any** of the negative matcher values from each provided category.
- **Must** match **any** of the positive matcher values from each provided category.

As a concrete example, consider an invocation like so:
```shell
… | krf --not-kind job,po --name backend,cache --not-label role=secondary --namespace client-a,client-b,client-c
```

The following checks are then performed against each resource in the input stream:

- If the resource is of kind of `Job`, then the resource is **rejected**.
- If the resource is of kind of `Pod`, then the resource is **rejected**.
- If the resource has the label `role=secondary`, then the resource is **rejected**.
- If the resource does not have a name of `backend` or `cache`, then the resource is **rejected**.
- If the resource does not have a namespace of `client-a`, `client-b`, or `client-c`, then the resource is **rejected**.
- Otherwise, the resource is finally **accepted**!

A resource such as this would be accepted by the above rules:
```yaml
apiVersion: apps/v1
kind: Deployment

metadata:
  name: cache
  namespace: client-b
  labels:
    role: primary
```

### Outputting Resources

By default, resources will be output in a table format:
```shell
… | krf
Namespace    API Version  Kind        Name              
─────────    ───────────  ────        ────              
default      apps/v1      Deployment  backend
```

If resources are sourced from a file (opposed to stdin) then the associated paths will also be included:
```shell
krf ./manifests
Namespace    API Version  Kind        Name     Path
─────────    ───────────  ────        ────     ────
default      apps/v1      Deployment  backend  ./manifests/deployment.yaml
```

If output is being redirected (to a file or another command) then the resources will be output in a yaml format:
```shell
… | krf | cat
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: backend
```

You can also explicitly choose to output resources in yaml format:
```shell
… | krf -o=yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  namespace: default
  name: backend
```

Conversely, you must explicitly choose to output resources in table format if you want to redirect that output: 
```shell
… | krf -o=table | tee summary.txt
Namespace    API Version  Kind        Name              
─────────    ───────────  ────        ────              
default      apps/v1      Deployment  backend
```

You can also output only the names of filtered resources (similar to `kubectl -o=name`):
```shell
… | krf -o=name
Deployment/backend
```

Or output only the names of resources that are referenced by filtered resources:

```shell
… | krf -o=references
ConfigMap/configuration
PersistentVolumeClaim/database
Secret/api-credentials
```

### Tips & Tricks

Here is a collection of some useful ways to utilize `krf`.

Get a quick overview of all resources:
```shell
krf ./kustomize
```

Show resource patches (files ending in `.patch.yaml`) in the production overlay:
```shell
krf ./kustomize/environments/production --patch
```

Identify deployments that would roll if the credentials ConfigMap (which may have a generated name like `credentials-8mbdf7882g`) were updated: 
```shell
kubectl get deploy -o=yaml | krf --references cm/credentials- 
````

Identify pods that do not have a security context configured:
```shell
kubectl get pod -o=yaml | krf --not-jsonpath '..securityContext'
```

Iteratively verify the impact of refactors on resources:
```shell
# Render out the current resources:
kustomize build ./kustomize/environments/production > original.yaml

# Perform your refactors on the various resources or manifests:
vim ./kustomize/applications/kustomization.yaml

# Render out the (now refactored) resources:
kustomize build ./kustomize/environments/production > modified.yaml

# Display a list of any of the modified resources which differences from their 
# original counterparts. Verify that your changes had the desired impact, or 
# even to verify that no resources were ultimately impacted.
krf --diff original.yaml modified.yaml
```

Selectively apply certain resources:
```shell
kustomize build … | krf --kind svc,ing | kubectl apply -f -
```

Perform some additional processing on resources:
```shell
krf … -o=json | while IFS= read -r line; do
  echo $line | jq -r '. | "processing \(.kind)/\(.metadata.name)..."'
done
```

## License

This code is distributed under the [MIT License][license-link], see [LICENSE.txt][license-file] for more information.

---

<p align="center">
  Created by <a href="https://github.com/joshdk">Josh Komoroske</a> ☕
</p>

[github-actions-badge]:  https://github.com/joshdk/krf/actions/workflows/build.yaml/badge.svg
[github-actions-link]:   https://github.com/joshdk/krf/actions/workflows/build.yaml
[github-release-badge]:  https://img.shields.io/github/release/joshdk/krf/all.svg
[github-release-link]:   https://github.com/joshdk/krf/releases
[license-badge]:         https://img.shields.io/badge/license-MIT-green.svg
[license-file]:          https://github.com/joshdk/krf/blob/master/LICENSE.txt
[license-link]:          https://opensource.org/licenses/MIT
