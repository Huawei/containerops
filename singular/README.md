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
kubernetes:
  version: 1.6.6 # or git commit hash tag: b516a521b14d86e59c5bb88930d20502a0712d78
  components:
    kubectl:
      url: https://hub.opshub.sh/containerops/singular/binary/kubectl/1.6.6
    kubelete:
      url: https://hub.opshub.sh/containerops/singular/binary/kubelete/1.6.6
  masters:
    number: 1
  nodes:
    number: 3
etcd:
  version: 3.1.1
  components:
    etcd:
      url: https://hub.opshub.sh/containerops/singular/binary/etcd/3.1.1
    etcdctl:
      url: https://hub.opshub.sh/containerops/singular/binary/etcdctl/3.1.1
  nodes:
    number: 3
```