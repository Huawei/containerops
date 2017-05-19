# Singular

### NAME

```
$ singular -- The deployment and operations tools.
```
### SYNOPSIS
```
Usage:	singular <command> <subcommand> [option] 
```
```
Available Commands:
init     Initialize your singular application config for installing kubernetes cluster
config   Config  Configure your nodes of kubernetes cluster
install  To start a new kubernetes cluster installing and running each services
cluster  List kubernetes cluster information and status
```

### DESCRIPTION SUBCOMMAND & OPTION
    
```
command:
init    Installing kubernetes cluster with default configration
option:
		--apikey	APIkey you have generated to access the public cloud API.
		--cerkey	Generated key-certificate pairs could help to access to the linux server    without the need to type password.
					by using Generated key-certificate pairs, you could access to the linux server without typing password.
		--cerpath	Without CApath option ,the default value is /etc/.singular/id_rsa.pub
					Or you could type your custom path for generate file id_rsa and id_rsa.pub
```
```
command:
config  Configure your nodes of kubernetes cluster
option:
		--master|slave                       Create master or slave nodes
		--security                           Generate kubernetes certificate
		--privtenet       					 Privte network for your cluster
		--count =3   <value>		    	 Number of nodes in cluster
        --mSize =512	 <1024|2048|>        Node memory Size
        --region =sfo    sfo|nyc			 Cluster's localization of the region
        --slug=ubuntu-17-04-x64  <value>     System version
```
```
command:
install Start to install kubenetes cluster automatically by the configuration file.
        --master|slave     Custom installation, master or slave nodes only.
        --pull             Download Kubernetes binaries without install.               
option:
        --config=~/etc/.singular/config.yaml      setting custom singular path of config.
```
```
command:
		cluster List kubernetes cluster information and status
option:
		--master|slave     List master or slave nodes information and status          
```
## Using singulary with a configuration file

##### It’s possible to configure singulary with a configuration file instead of command line flags, and some more advanced features may only be available as configuration file options. 

###Sample  Configuration

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