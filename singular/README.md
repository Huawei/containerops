## Singular

Singular is designed to deploy ContainerOps platform, mainly focus on Kubernetes, Prometheus, and other Cloud Native technology stack. The goal is to build the stack across clouds, OpenStack, VMs and even bare metals.

Singular doesn't use any other deploy tools like _kubeadm_, but deploys everything in a hard way instead. The **hard** way means, user prepare all the binaries, nodes and SSH keys to nodes, then singular will copy the binaries to the nodes, generate configs, CA, systemd scripts, and start the service. The process is just like you deploy a kubernetes cluster manually.

Singular provides templates of different service version combinations, currently the services include [etcd](https://github.com/coreos/etcd), [flannel](https://github.com/coreos/flannel), [docker-ce](https://github.com/docker/docker-ce) and [Kubernetes](https://github.com/kubernetes/kubernetes). More CNCF projects will be supported in the future.

### Deployment Template

Singular uses a **YAML** file as deployment template, it describes the architecture of the cluster, including the nodes' info, the services' metadata like binary path, version etc. The sample templates can be found in `./external/public/data/cncf.build`. The structure is like this:

```yaml
uri: containerops/singular/etcd-3.2.8-flanneld-0.7.1-docker-17.04.0-ce-k8s-1.9.2
title: Deploy Kubernetes With etcd-3.2.8 flanneld-0.7.1 docker-17.04.0-ce k8s-1.9.2
tag: latest
nodes: # 3 nodes is provided at least, the first one will be deployed as kubernetes master node
  -
    ip: 192.168.43.70
    user: root
    distro: archlinux # distro could be 'ubuntu', 'centos' or 'archlinux'
  -
    ip: 192.168.43.71
    user: root
    distro: archlinux
  -
    ip: 192.168.43.72
    user: root
    distro: archlinux
tools:
  ssh: # The ssh key to target nodes
    private: /home/lance/.ssh/id_rsa
    public: /home/lance/.ssh/id_rsa.pub
infras:
  -
    name: etcd
    version: etcd-3.2.8
    master: 3
    minion: 0
    components:
      -
        binary: etcd
        url: /my/k8s-binaries/etcd/3.2.8/etcd  # The value can be file path or http url
        package: false
        systemd: etcd-3.2.8
        ca: etcd-3.2.8
      -
        binary: etcdctl
        url: /my/k8s-binaries/etcd/3.2.8/etcdctl
        package: false
  -
......
```


#### SSH Key

When deploy infrastructures, **Singular** need to _SSH_ to virtual machines or bare metals.

1. At least provide SSH private key, **Singular** will create public key from it. Generate ssh key follow Github document - [How to generate SSH Key](https://help.github.com/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent).
2. If no SSH key files in deployment template, **Singular** will create _SSH_ pair key files in default folder(**_$HOME/.containerops/ssh_**) and name(**_id_rsa.pub_** and **_id_rsa_**).

#### Deploy Command

```
singular deploy template /tmp/deploy.yml  --verbose --timestamp
```
#### The next step

After deploying the cluster, the in-cluster DNS is not installed by default, a kubernetes cluster usually use CoreDNS or KubeDNS as name server. The yaml files are already there for the newly deployed cluster, you can set them up by:
```
kubectl apply -f $GOPATH/src/github/Huawei/containerops/singular/external/public/data/dns/kubedns.yaml
```
Or CoreDNS:
```
kubectl apply -f $GOPATH/src/github/Huawei/containerops/singular/external/public/data/dns/coredns.yaml
```

The `__PILLAR__DNS__SERVER__` and `__PILLAR__DNS__DOMAIN__` are preconfigured according to the kubelet template file(that it, 10.254.0.2 and cluster.local.). Feel free to change them to meet your own requirements.
