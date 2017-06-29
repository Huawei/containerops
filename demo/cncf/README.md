## ContainerOps Demo With CNCF Technology Stack

### The Demo Design

![CNCF Demo Design](demo-for-cncf-ci.png)

Quinton Hoole ([@quintonhoole](https://github.com/quintonhoole)) and Quanyi Ma ([@genedna](https://github.com/genedna)) had a [proposal](https://docs.google.com/document/d/1G2UXaDBjGXpdvVD-yl-hpIV4z8RRJvEvaJhsPEBbiPI) for [CNCF CI Working Group](https://github.com/cncf/wg-ci). I designed a demo followed the proposal, then implemented with ContainerOps components and orchestration engine.

### How to try the demo?

1. Should have a Kubernetes cluster, I suggest use the Google Container Engine.
2. Should have `kubectl` in the local system, and it connecting to cluster successfully.
3. Compile the `pilotage` command in the `pilotage` folder.
4. Run with command:

```
pilotage cli run cncf-demo.yaml --verbose true --timestamp true
```

### What is the demo doing?

There are two phases of the demo: 
1. Build and test the three projects: Kubernetes, Prometheus, CoreDNS, then compile their binaries.
2. Build a Native Cloud stack with binaries in the DigitalOcean cloud. Then deploy the CNCF demo in the 

### Working in process

1. The build/test/release components for Kubernetes, Prometheus and CoreDNS almost done. We are working on deploy the Dockyard which is a container registry and binary repository.
2. Also we are working on deploy Kubernetes/Prometheus/CoreDNS/Etcd/Flannel in DigitalOcean.
3. We choose the CNCF official [demo](https://github.com/cncf/demo) to deploy.

### CNCF Demo Configuration File [Status: _WIP_ ]

```yaml
name: cncf/demo-for-cncf-ci
title: Demo For Cloud Native Computing Foundation CI Working Group
version: 4
tag: latest
timeout: 0
stages:
  -
    type: start
    name: start
    title: Start
  -
    type: normal
    name: prometheus-test-build-release
    title: Building, testing Prometheus project, compile then upload to Dockyard artifact repository.
    sequencing: sequence
    actions:
      -
        name: build-prometheus
        title: Build Prometheus project with "make build"
        jobs:
          -
            type: component
            kubectl: prometheus/prometheus-build.yaml
            endpoint: docker.io/containerops/cncf-demo-prometheus:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "prometheus=https://github.com/prometheus/prometheus.git action=build release=test.opshub.sh/containerops/cncf-demo/demo"
            output: []
      -
        name: test-prometheus
        title: Test Prometheus project with "make test"
        jobs:
          -
            type: component
            kubectl: prometheus/prometheus-test.yaml
            endpoint: docker.io/containerops/cncf-demo-prometheus:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "prometheus=https://github.com/prometheus/prometheus.git action=test release=test.opshub.sh/containerops/cncf-demo/demo"
            output: []
      -
        name: release-prometheus
        title: Compile Prometheus project with "make build", then upload to artifact repository
        jobs:
          -
            type: component
            kubectl: prometheus/prometheus-release.yaml
            endpoint: docker.io/containerops/cncf-demo-prometheus:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "prometheus=https://github.com/prometheus/prometheus.git action=release release=test.opshub.sh/containerops/cncf-demo/demo"
            output: ["CO_PROMETHEUS_URI", "CO_PROMTOOL_URI"]
  -
    type: normal
    name: coredns-test-build-release
    title: Building, testing CoreDNS project, compile then upload to Dockyard artifact repository.
    sequencing: sequence
    actions:
      -
        name: build-coredns
        title: Build CoreDNS project with "make coredns"
        jobs:
          -
            type: component
            kubectl: coredns/coredns-build.yaml
            endpoint: docker.io/containerops/cncf-demo-coredns:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "coredns=https://github.com/coredns/coredns.git action=test release=test.opshub.sh/containerops/cncf-demo/demo"
            output: []
      -
        name: test-coredns
        title: Test CoreDNS project with "make test"
        jobs:
          -
            type: component
            kubectl: coredns/coredns-test.yaml
            endpoint: docker.io/containerops/cncf-demo-coredns:latest
            resources:
              cpu: 2
              memory: 4G
            timeout: 0
            environments:
              - CO_DATA: "coredns=https://github.com/coredns/coredns.git action=test release=test.opshub.sh/containerops/cncf-demo/demo"
            output: []
      -
        name: release-coredns
        title: Compile CoreDNS project with "make coredns", then upload to artifact repository
        jobs:
          -
            type: component
            kubectl: coredns/coredns-release.yaml
            endpoint: docker.io/containerops/cncf-demo-coredns:latest
            resources:
              cpu: 0
              memory: 0
            timeout: 0
            environments:
              - CO_DATA: "coredns=https://github.com/coredns/coredns.git action=release release=test.opshub.sh/containerops/cncf-demo/demo"
            output: ["CO_COREDNS_URI"]
  -
    type: normal
    name: k8s-test-build-release
    title: Building, testing Kubernetes project using Bazel, compile then upload to Dockyard artifact repository.
    sequencing: sequence
    actions:
      -
        name: build-kubernetes
        title: Build Kubernetes project with "make bazel-build"
        jobs:
          -
            type: component
            kubectl: kubernetes/kubernetes-build.yaml
            endpoint: docker.io/containerops/cncf-demo-kubernetes:latest
            resources:
              cpu: 2
              memory: 8G
            timeout: 0
            environments:
              - CO_DATA: "kubernetes=https://github.com/kubernetes/kubernetes.git action=build release=test.opshub.sh/containerops/cncf-demo/demo"
            output: []
      -
        name: release-kubernetes
        title: Compile Kubernetes project with "make all", then upload all binaries to artifact repository
        jobs:
          -
            type: component
            kubectl: kubernetes/kubernetes-release.yaml
            endpoint: docker.io/containerops/cncf-demo-kubernetes:latest
            resources:
              cpu: 2
              memory: 8G
            timeout: 0
            environments:
              - CO_DATA: "kubernetes=https://github.com/kubernetes/kubernetes.git action=release release=test.opshub.sh/containerops/cncf-demo/demo"
            output: ["CO_APIEXTENSIONS-APISERVER_URI", "CO_CLOUD-CONTROLLER-MANAGER_URI", "CO_CONVERSION-GEN_URI", "CO_DEEPCOPY-GEN_URI", "CO_DEFAULTER-GEN_URI", "CO_E2E.TEST_URI", "CO_E2E_NODE.TEST_URI", "CO_GENDOCS_URI", "CO_GENFEDDOCS_URI", "CO_GENKUBEDOCS_URI", "CO_GENMAN_URI", "CO_GENSWAGGERTYPEDOCS_URI", "CO_GENYAML_URI", "CO_GINKGO_URI", "CO_GKE-CERTIFICATES-CONTROLLER_URI", "CO_GO-BINDATA_URI", "CO_HYPERKUBE_URI", "CO_KUBE-AGGREGATOR_URI", "CO_KUBE-APISERVER_URI", "CO_KUBE-CONTROLLER-MANAGER_URI", "CO_KUBE-PROXY_URI", "CO_KUBE-SCHEDULER_URI", "CO_KUBEADM_URI", "CO_KUBECTL_URI", "CO_KUBEFED_URI", "CO_KUBELET_URI", "CO_KUBEMARK_URI", "CO_LINKCHECK_URI", "CO_MUNGEDOCS_URI", "CO_OPENAPI-GEN_URI", "CO_TESTSTALE_URI"]
      -
        name: test-kubernetes
        title: Test Kubernetes project with "make bazel-test"
        jobs:
          -
            type: component
            kubectl: kubernetes/kubernetes-test.yaml
            endpoint: docker.io/containerops/cncf-demo-kubernetes:latest
            resources:
              cpu: 2
              memory: 8G
            timeout: 0
            environments:
              - CO_DATA: "kubernetes=https://github.com/kubernetes/kubernetes.git action=test release=test.opshub.sh/containerops/cncf-demo/demo"
            output: []
  -
    type: end
    name: end
    title: End

```