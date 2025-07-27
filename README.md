# Multitenancy Operator

## Overview

This is a simple proof-of-concept Kubernetes operator designed to demonstrate how custom resources can drive automation. It watches a custom resource called `TenantInfo`, and for each defined tenant, it creates a corresponding `ConfigMap` containing tenant-specific data. While this operator currently only manages `ConfigMaps`, these are used as placeholders to illustrate how other more complex resources could be managed in a real-world scenario.

Some examples of what a production-grade operator might automate include:

- Deploying an application on demand
- Taking and restoring backups of that application's state
- Handling upgrades of the application code alongside related changes such as database schemas or extra configuration settings
- Publishing a Service to applications that don't support Kubernetes APIs to discover them
- Simulating failure in all or part of your cluster to test its resilience
- Choosing a leader for a distributed application without an internal member election process

## Description

This operator manages multitenant configurations using a Kubernetes-native approach. It introduces a custom resource definition (CRD) named `TenantInfo`, where each resource can define a list of tenants and their associated `webserviceURL`. The operator watches these CRs and ensures that a `ConfigMap` is created or updated for each tenant, reflecting the declared state.

Designed as a lightweight and extensible example, this project demonstrates the reconciliation pattern and how it can be used to bridge declarative specifications with automated infrastructure behavior. Although it currently focuses on `ConfigMaps`, the architecture is ready to be extended to manage Deployments, Services, Secrets, database provisioning, routing rules, or any other resources that need to be tenant-aware.

This project demonstrates how to build custom operators using Kubebuilder and the controller-runtime framework.

## Setup

This section initializes the project, generates boilerplate code, and installs the custom resource definitions (CRDs) into the cluster. It's required only once when bootstrapping or resetting the operator project.

```shell
# Initialize the operator project with your domain and module path
kubebuilder init --domain example.com --repo github.com/mustafa-qamaruddin/multitenancy-operator

# Create a new API with group/version/kind -> TenantInfo CRD
kubebuilder create api --group multitenancy-management --version v1 --kind TenantInfo

# Generate Go code for the API types and deepcopy functions
make generate

# Generate CRD YAMLs and RBAC permissions from annotations
make manifests

# Apply the generated CRDs and RBAC roles to your cluster
make install
```

## Run Demo

This section runs the operator locally and applies a sample `TenantInfo` custom resource. You can observe the operator creating `ConfigMaps` in response. The dashboard is optional but useful for visual inspection.

```shell
# Open the Kubernetes dashboard (optional but helpful for debugging)
minikube dashboard

# Run the operator controller locally, watching your cluster for CR changes
make run

# Apply the sample TenantInfo CR to trigger reconciliation logic
kubectl apply -k config/samples

# List all TenantInfo resources across all namespaces
kubectl get tenantinfoes -A

# Delete the sample to test if associated resources (like ConfigMaps) are cleaned up
kubectl delete tenantinfoes/tenantinfo-sample
```

## Cleanup

This section uninstalls all the CRDs and RBAC settings that were installed during the setup. Run this to clean your cluster after testing.

```shell
# Remove all CRDs and associated RBAC resources created by the operator
make uninstall
```

## Tests

This section shows how to run the automated tests 

```shell
# To ensure the binaries are properly set up
make setup-envtest

# For e2e test create a node in the cluster
brew install kind
kind create cluster --name kind

# Runs all global tests
go test ./...

# Runs only the controller tests
go test ./internal/controller/...
```

## KubeBuilder Default README.md

## Getting Started

### Prerequisites
- go version v1.24.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster
**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/multitenancy-operator:tag
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

```sh
make install
```

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/multitenancy-operator:tag
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Create instances of your solution**
You can apply the samples (examples) from the config/sample:

```sh
kubectl apply -k config/samples/
```

>**NOTE**: Ensure that the samples has default values to test it out.

### To Uninstall
**Delete the instances (CRs) from the cluster:**

```sh
kubectl delete -k config/samples/
```

**Delete the APIs(CRDs) from the cluster:**

```sh
make uninstall
```

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following the options to release and provide this solution to the users.

### By providing a bundle with all YAML files

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/multitenancy-operator:tag
```

**NOTE:** The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without its
dependencies.

2. Using the installer

Users can just run 'kubectl apply -f <URL for YAML BUNDLE>' to install
the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/multitenancy-operator/<tag or branch>/dist/install.yaml
```

### By providing a Helm Chart

1. Build the chart using the optional helm plugin

```sh
kubebuilder edit --plugins=helm/v1-alpha
```

2. See that a chart was generated under 'dist/chart', and users
can obtain this solution from there.

**NOTE:** If you change the project, you need to update the Helm Chart
using the same command above to sync the latest changes. Furthermore,
if you create webhooks, you need to use the above command with
the '--force' flag and manually ensure that any custom configuration
previously added to 'dist/chart/values.yaml' or 'dist/chart/manager/manager.yaml'
is manually re-applied afterwards.

## Contributing
// TODO(user): Add detailed information on how you would like others to contribute to this project

**NOTE:** Run `make help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

