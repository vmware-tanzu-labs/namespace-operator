# namespace-operator

This project was built using [operator-builder](https://github.com/vmware-tanzu-labs/operator-builder).  Please visit
the link for more information on how to build your operators for management of Kubernetes workloads.

## Project Motivation

The namespace-operator project is a project designed for implementing best practices for each namespace created
within a Kubernetes environment.  More specifically, it was implemented to take into account best practices
from several years of production-level, field experience from the engineers currently working at VMware.  While it is
noted that there are several challenges with multi-tenant clusters, for those who desire multi-tenancy within cluster can
use the namespace-operator model to implement a level of tenancy within each Kubernetes cluster which satisfies many production
requirements.  Best practices for namespace-operator are derived from guides in the Tanzu Developer Center
(Workload Tenancy Guide) as documented at https://tanzu.vmware.com/developer/guides/kubernetes/workload-tenancy/.

## Architecture Overview

The architecture for namespace-operator is built off of a single construct called a `TanzuNamespace`. This construct
is a Kubernetes `CustomResourceDefinition` which is responsible for the creation and reconciliation of several other
subordinate constructs such as the following:

- `Namespace` - first and foremost, a namespace is created based off of each `TanzuNamespace` object.  The namespace
  is in place to represent a tenant for a workload.
- `LimitRange` - for each `TanzuNamespace`, a limit range is created to give a sane resource range for workloads that reside
  within the namespace (tenant).  See https://kubernetes.io/docs/concepts/policy/limit-range/.  Limit ranges may be adjusted
  but are defaulted to a low value by default to encourage administrators to specify the values that they need.
- `ResourceQuota` - for each `TanzuNamespace`, a resource quota is created to provide a sane limitation for resources to
  be consumed within that namespace.  See https://kubernetes.io/docs/concepts/policy/resource-quotas/.  Resource
  quotas may be adjusted but are defaulted to a low value by default to encourage administrators to specify the values that
  they need.
- `NetworkPolicy` - for each `TanzuNamespace`, a network policy is created to provide microsegmentation for the namespace/tenant.
  See https://kubernetes.io/docs/concepts/services-networking/network-policies/.  In the case of a `TanzuNamespace`, we implement
  a default `deny-all` policy and only allow traffic out for DNS queries.  This provides a namespace lockdown by default and
  forces users to define which ports via their own `NetworkPolicy` truly should be allowed.  **NOTE:** network policy implementation is
  highly dependent on the Kubernetes CNI selection.  Please ensure your CNI implements the NetworkPolicy spec to use.
- **RBAC (Not Yet Implemented)** - for each `TanzuNamespace`, the namespace-operator will lay down some role-based access
  control to implement the workload administrator, developer, and view roles along with associated service accounts and role bindings.
  **NOTE:** this is not yet implemented as of this writing.
- **ImagePullSecret (Not Yet Implemented)** - for each `TanzuNamespace`, an `ImagePullSecret` is created to allow workloads
  in the namespace to pull images from private image repositories.

## Architecture Diagram

![namespace-operator diagram](img/namespace-operator.png "namespace-operator diagram")

**NOTE:** RBAC and ImagePullSecret are not yet implemented.  Diagram to represent idea only.

## Installation

Run the following commands to install the namespace-operator:

1. Install the CRDs:

```bash
make install
```

2. Deploy the Docker Image:

```bash
IMG=ghcr.io/vmware-tanzu-labs/namespace-operator:v0.2.0 make deploy
```

3. (Optional) Install the Sample CRD as a test:

```bash
kubectl apply -f config/samples/tenancy_v1alpha2_tanzunamespace.yaml
```

The above commands will install the following objects into your cluster:

- TanzuNamespace CustomResourceDefinition
- RBAC for namespace-operator deployment
- namespace-operator deployment
- (Optional) Sample Namespace, LimitRange, ResourceQuota, NetworkPolicy

## Usage

The following is a representation of a `TanzuNamespace` resource definition (see `config/samples/tenancy_v1alpha2_tanzunamespace.yaml`):

```yaml
---
apiVersion: tenancy.platform.cnr.vmware.com/v1alpha2
kind: TanzuNamespace
metadata:
  name: tanzunamespace-sample
spec:
  namespace: "tanzu-namespace"
  resources:
    limits:
      cpu: "250m"
      memory: "256Mi"
    requests:
      cpu: "250m"
      memory: "256Mi"
    max:
      cpu: "500m"
      memory: "256Mi"
    quota:
      requests:
        cpu: "2000m"
        memory: "4Gi"
      limits:
        cpu: "2000m"
        memory: "4Gi"
```

The above can be applied via standard `kubectl apply -f <tanzu_namespace_file>`, substituting the appropriate values as necessary.
