## Singular

Singular design for deploy and operation ContainerOps platform, mostly focus on Kubernetes, Prometheus, others in the Cloud Native technology stack. We are trying to build all stack cross cloud, OpenStack, bare metals.

Singular don't use any other deploy tools like _kubeadm_, deploy everything in a hard way instead. Singular providers templates for service so could deploy any version. Singular deploys development versions of **Kubernetes**, **CoreDNS**, others in **CNCF** CI demo.

### Deployment Template

Singular uses a **YAML** file as deploy template, it describes the architecture of the cluster. It don't only deploy ContainerOps modules, also use to deploy Kubernetes, others in Cloud Native stack and some common software.

#### Template Samples

```YAML
uri: containerops/demo-for-cncf-ci/deploy-cncf-stack
title: Demo For Deploy Cloud Native Computing Foundation CI Working Group
version: 4
tag: latest
service:
  provider: digitalocean
  token: b516a521b14d86e59c5bb8893
tools:
  ca-tools:
    component:
      cfssl:
        url: https://hub.opshub.sh/containerops/singular/binary/cfssl/latest
      cfssl-certinfo:
        url: https://hub.opshub.sh/containerops/singular/binary/cfssl-certinfo/latest
      cfssljson:
        url: https://hub.opshub.sh/containerops/singular/binary/cfssljson/latest
infra:
  etcd:
    version: 3.2.2
    components:
      etcd:
        url: https://hub.opshub.sh/containerops/singular/binary/etcd/3.2.2
      etcdctl:
        url: https://hub.opshub.sh/containerops/singular/binary/etcdctl/3.2.2
    nodes: 3
    distro: false
  flannel:
    version: 0.7.1
    components:
      flanneld:
        url: https://hub.opshub.sh/containerops/singular/binary/flanneld/0.7.1
      scripts:
        url: https://hub.opshub.sh/containerops/singular/binary/mk-docker-opts.sh/0.7.1
    nodes: 0
    distro: false
  kubernetes:
    version: 1.6.7 # or git commit hash tag: b516a521b14d86e59c5bb88930d20502a0712d78
    components:
      cloud-controller-manager:
        url: https://hub.opshub.sh/containerops/singular/binary/cloud-controller-manager/1.6.7
      hyperkube:
        url: https://hub.opshub.sh/containerops/singular/binary/hyperkube/1.6.7
      kubeadm:
        url: https://hub.opshub.sh/containerops/singular/binary/kubeadm/1.6.7
      kube-aggregator:
        url: https://hub.opshub.sh/containerops/singular/binary/kube-aggregator/1.6.7
      kube-apiserver:
        url: https://hub.opshub.sh/containerops/singular/binary/kube-apiserver/1.6.7
      kube-controller-manager:
        url: https://hub.opshub.sh/containerops/singular/binary/kube-controller-manager/1.6.7
      kubectl:
        url: https://hub.opshub.sh/containerops/singular/binary/kubectl/1.6.7
      kubefed:
        url: https://hub.opshub.sh/containerops/singular/binary/kubefed/1.6.7
      kubelet:
        url: https://hub.opshub.sh/containerops/singular/binary/kubelet/1.6.7
      kube-proxy:
        url: https://hub.opshub.sh/containerops/singular/binary/kube-proxy/1.6.7
      kube-scheduler:
        url: https://hub.opshub.sh/containerops/singular/binary/kube-scheduler/1.6.7                 
    nodes:
      master: 1
      node: 3
    distro: false
  docker:
    version: 1.7.04.0-ce
    components:
      docker:
        url: https://hub.opshub.sh/containerops/singular/binary/docker/17.04.0-ce
      dockerd:
        url: https://hub.opshub.sh/containerops/singular/binary/dockerd/17.04.0-ce
      docker-init:
        url: https://hub.opshub.sh/containerops/singular/binary/docker-init/17.04.0-ce
      docker-proxy:
        url: https://hub.opshub.sh/containerops/singular/binary/docker-proxy/17.04.0-ce
      docker-runc:
        url: https://hub.opshub.sh/containerops/singular/binary/docker-runc/17.04.0-ce                                     
      docker-containerd:
        url: https://hub.opshub.sh/containerops/singular/binary/docker-containerd/17.04.0-ce
      docker-containerd-ctr:
        url: https://hub.opshub.sh/containerops/singular/binary/docker-containerd-ctr/17.04.0-ce
      docker-containerd-shim:
        url: https://hub.opshub.sh/containerops/singular/binary/docker-containerd-shim/17.04.0-ce
    nodes: 0
    distro: false
```