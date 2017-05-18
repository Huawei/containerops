# Singular

## NAME

```
     singular -- The deployment and operations tools.
```
## SYNOPSIS
```
    Usage:	singular [OPTIONS] [ARG...]
```
```
    The following options are available:
    [--user <name> ]
    [--tokens <value>] 
    [--ssh] <directory_name>
    [--cluster [master|node]][--NodeCount][--MSize][--Region][--Slug][--Privtenet]
    [--security yes|no]
    [--install master|node]
```

## DESCRIPTION
    
```
    --user  Your custom name or system account. default is "singular_user"
```
```
    --tokens    Tokens you have generated to access the vm cloud API.
```
```
    --ssh   Setting up SSH keys that access to your Linux server without the need for type password.
            Without ssh option ,the default value is /etc/singular/id_rsa.pub
            Or you could type your custom path for generate file id_rsa and id_rsa.pub
```
```  					
    --security  singular will automatically generate kubernetes certificate for provides an additional layer of security. 
```
```
    --cluster    
            master|node  <NodeCount>            Custom configurations for master|node vm 
            NodeCount   <value>				 Number of nodes in cluster
            MSize		<512|1024|2048|>         Node memory Size
            Region    sfo|nyc					 Node local region
            Slug      <value>					 System version
            Privtenet yes|no					 Privte network for your cluster
```
```
    --install Start to install kubenetes cluster automatically by the configuration file. "/etc/singular/config.yaml"
```
# Using singulary with a configuration file

### It’s possible to configure singulary with a configuration file instead of command line flags, and some more advanced features may only be available as configuration file options. 

Sample  Configuration
```
cluster_config:
    User:     "singular_user"
    Token:     ""
    EtcdNet: "/kube/network"
    Security: "yes"

vm_config:
    MSize:     "1024mb"
    Region:     "sfo2"
    Slug:     "ubuntu-17-04-x64"
    Fingerprint:      "ee:81:d0:59:ab:09:1c:ff:52:dd:11:f8:bd:a6:7f:a8"
cluster_download:
    cloud_controller_manager : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/cloud-controller-manager"
    hyperkube : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/hyperkube"
    kube_aggregator : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-aggregator"
    kube_apiserver : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-apiserver"
    kube_controller_manager : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-controller-manager"
    kube_proxy : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-proxy"
    kube_scheduler : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-scheduler"
    kubeadm : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kubeadm"
    kubectl : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kubectl"
    kubefed : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kubefed"
    kubelet : "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kubelet"
    
    
```