# Kubernetes Admission Webhook Server

[![Version](https://img.shields.io/github/v/release/52North/admission-webhook-server)](https://github.com/52North/admission-webhook-server/releases)
[![GoDoc](https://godoc.org/github.com/52North/admission-webhook-server?status.svg)](https://godoc.org/github.com/52North/admission-webhook-server)
![Workflow](https://github.com/52North/admission-webhook-server/workflows/Release/badge.svg)

---

API server providing webhook endpoints for Kubernetes admission controller to mutate objects.

Currently it can handle mutating `nodeSelector` based on namespaces. This same functionality exists in standard Kubernetes cluster installation if enabled. However it's not enabled in EKS.

The server can be easily extended by adding more handlers for different mutations needs.

The repo also includes a Helm chart for easy deployment to your Kubernetes cluster.

---

## Installation

You need to update helm value `podNodesSelectorConfig` in `chart/values.yaml` so it can use the value to mutate the pods.

Note: below example using Helm v3. However the chart is compatible with helm version older than v3.

```sh
$ git clone https://github.com/52North/admission-webhook-server
$ cd admission-webhook-server/helm
$ helm install admission-webhook-server .
```

## Helm

The following table lists the configuration parameters for the helm chart.

| Parameter  | Description  | Default  |
|---|---|---|
| nameOverride  | Override general resource name   |   |
| basePathOverride  | Url base path   | mutate  |
| podNodesSelectorConfig  | Configuration for podnodesselector. The namespace and labels are set here following the format: namespace: key=label,key=label; namespace2: key=label. Multiple namespaces separate by ";". Example: develop: node-role.kubernetes.io/development=true, beta.kubernetes.io/instance-type=t3.large  |   |
| podTolerationRestrictionConfig | Configuration for podtolerationrestriction, a JSON object mapping namespaces to lists of tolerations: `{"namespace": [{"operator": "Equal", "effect": "NoSchedule", "key": "some-taint", "value": "some-taint-value"}]}'` | |
| service.name  | Name of the service. It forms part of the ssl CN  | admission-webhook  |
| service.annotations  | Annotation for the service  | {} |
| replicas | Number of replicas  | 1  |
| strategy.type  | Type of update strategy  | RollingUpdate  |
| image  | Docker image name  | 52north/admission-webhook-server  |
| imageTag  | Docker image tag  | latest  |
| imagePullPolicy  | Docker image pull policy  | Always  |
