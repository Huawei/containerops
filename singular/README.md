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
nodes: 3
service:
  provider: digitalocean
  token: b516a521b14d86e59c5bb8893
  region: sfo2
  size: 4G
  image: ubuntu-17-04-x64
tools:
  ssh:
    id-rsa: $HOME/.containerops/ssh/id-rsa
    id-rsa.pub: $HOME/.containerops/ssh/id-rsa.pub
infra:
  etcd:
    version: 3.2.2
    components:
      etcd:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/etcd/3.2.2
        package: false
        systemd: etcd-3.2.2
        ca: etcd-3.2.2
      etcdctl:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/etcdctl/3.2.2
        package: false
    nodes: 3
  flannel:
    version: 0.7.1
    components:
      flanneld:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/flanneld/0.7.1
        package: false
        systemd: flannel-0.7.1
        ca: flannel-0.7.1
        before: "etcdctl --endpoints={{.EtcdEndpoints}} --ca-file={{.CaPemFile}} --cert-file={{.FlanneldPemFile}} --key-file={{.FlanneldKeyFile}}"
      scripts:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/mk-docker-opts.sh/0.7.1
        package: false
    nodes: 0
  kubernetes:
    version: 1.6.7
    components:
      kube-apiserver:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/kube-apiserver/1.6.7
        package: false
        systemd: kube-apiserver-1.6.7
        ca: kubernetes-1.6.7
      kube-controller-manager:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/kube-controller-manager/1.6.7
        package: false
        systemd: kube-controller-manager-1.6.7
        ca: kubernetes-1.6.7
      kube-scheduler:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/kube-scheduler/1.6.7
        package: false
        systemd: kube-scheduler-1.6.7
      kubectl:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/kubectl/1.6.7
        package: false
      kubelet:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/kubelet/1.6.7
        package: false
        systemd: kubelet-1.6.7
      kube-proxy:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/kube-proxy/1.6.7
        package: false
        systemd: kube-proxy-1.6.7
        ca: kube-proxy-1.6.7
    nodes:
      master: 1
      node: 3
    package: false
  docker:
    version: 1.7.04.0-ce
    components:
      docker:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/docker/17.04.0-ce
        package: false
        systemd: docker-1.7.04.0-ce
        before: "apt update && apt dist-upgrade && apt install -y bridge-utils aufs-tools cgroupfs-mount libltdl7 && systemctl stop ufw && systemctl disable ufw && iptables -F && iptables -X && iptables -F -t nat && iptables -X -t nat"
        after: "iptables -P FORWARD ACCEPT"
      dockerd:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/dockerd/17.04.0-ce
        package: false
      docker-init:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/docker-init/17.04.0-ce
        package: false
      docker-proxy:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/docker-proxy/17.04.0-ce
        package: false
      docker-runc:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/docker-runc/17.04.0-ce
        package: false
      docker-containerd:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/docker-containerd/17.04.0-ce
        package: false
      docker-containerd-ctr:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/docker-containerd-ctr/17.04.0-ce
        package: false
      docker-containerd-shim:
        url: https://hub.opshub.sh/binary/v1/containerops/singular/binary/docker-containerd-shim/17.04.0-ce
        package: false
    nodes: 0
```