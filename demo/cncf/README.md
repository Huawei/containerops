## ContainerOps Demo With CNCF Technology Stack

### The Demo Design

![CNCF Demo Design](demo-for-cncf-ci.png)

Quinton Hoole ([@quintonhoole](https://github.com/quintonhoole)) and Quanyi Ma ([@genedna](https://github.com/genedna)) had a [proposal](https://docs.google.com/document/d/1G2UXaDBjGXpdvVD-yl-hpIV4z8RRJvEvaJhsPEBbiPI) for [CNCF CI Working Group](https://github.com/cncf/wg-ci). I designed a demo followed the proposal, then implemented with ContainerOps components and orchestration engine.

### How to try the demo?

1. *_Repositories_* Any commit or pull request trigger the demo. If you try this demo, you need to fork the Kubernetes, Prometheus, CoreDNS first. Then you should update the git repository URI in the cncf-demo YAML file.
2. *_Cloud_* The demo only supports DigitalOcean just now, so you need to registry an account of DigitalOcean. Then you should generate an access token for the demo and update the deploy component input data.
3. *_Debug_* Run the pilotage with the YAML file.

```
pilotage run --local-flow=cncf-demo.yaml
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
name: cncnf/demo-for-cncf-ci
title: Demo For Cloud Native Computing Foundation CI Working Group
version: 12
tag: latest
timeout: 0
stages: 
  - stage: 
      type: start
      name: start
      title: Start
  - stage:
      type: normal
      name: k8s-test-build-release
      title: Build, Test and Release the Kubernetes Binary
      sequencing: sequence # parallel or sequence
      actios: 
        - action:
            name: build-kubernetes
            title: Build Kubernetes
            jobs:
              - job:
                  type: component # component or service
                  endpoint: docker.io/containerops/cncf-demo-kubernetes
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: kubernetes=https://github.com/kubernetes/kubernetes.git action=build
        - action:    
            name: test-kubernetes
            title: Test Kubernetes
            jobs:
              - job:
                  type: component
                  endpoint: docker.io/containerops/cncf-demo-kubernetes
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: kubernetes=https://github.com/kubernetes/kubernetes.git action=test
        - action:    
            name: release-kubernetes
            title: Release Kubernetes
            jobs:
              - job:
                  type: component
                  endpoint: docker.io/containerops/cncf-demo-kubernetes
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: kubernetes=https://github.com/kubernetes/kubernetes.git action=release
  - stage:
      type: normal
      name: prometheus-test-build-release
      title: Build, Test and Release the Prometheus Binary
      sequencing: sequence
      actions:
        - action:
            name: build-prometheus
            title: Build Prometheus
            jobs:
              - job:
                  type: component # component or service
                  endpoint: docker.io/containerops/cncf-demo-prometheus
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: prometheus=https://github.com/prometheus/prometheus.git action=build
        - action:    
            name: test-prometheus
            title: Test Prometheus
            jobs:
              - job:
                  type: component
                  endpoint: docker.io/containerops/cncf-demo-prometheus
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: prometheus=https://github.com/prometheus/prometheus.git action=test
        - action:    
            name: release-prometheus
            title: Release Prometheus
            jobs:
              - job:
                  type: component
                  endpoint: docker.io/containerops/cncf-demo-prometheus
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: prometheus=https://github.com/prometheus/prometheus.git action=release
  - stage:
      type: normal
      name: coredns-test-build-release
      title: Build, Test and Release the CoreDNS Binary
      sequencing: sequence
      actions:
        - action:
            name: build-coredns
            title: Build CoreDNS
            jobs:
              - job:
                  type: component # component or service
                  endpoint: docker.io/containerops/cncf-demo-coredns
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: coredns=https://github.com/coredns/coredns.git action=build
        - action:    
            name: test-coredns
            title: Test CoreDNS
            jobs:
              - job:
                  type: component
                  endpoint: docker.io/containerops/cncf-demo-coredns
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: coredns=https://github.com/coredns/coredns.git action=test
        - action:    
            name: release-coredns
            title: Release CoreDNS
            jobs:
              - job:
                  type: component
                  endpoint: docker.io/containerops/cncf-demo-coredns
                  resources:
                    cpu: 0
                    memory: 0
                  timeout: 0 
                  environments:
                    - default: coredns=https://github.com/coredns/coredns.git action=release     
  - stage:
        type: end
        name: end
        title: End

```