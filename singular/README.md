# Singular

### Application


```
$ singular
Usage: singular [OPTIONS] COMMAND [arg...]
       singular [ --help | -v | --version ]

Singular, the Kubernetes deployment and operations tools.
```
####To automatically deploy Kubernetes, you could simply follow below steps:
### Before You Start
Before using the singular, you need to tell it about your public cloud credentials. Two steps as following:
 
##### 1）  Register an account for public cloud,  get API server key and set it into the Singular.
##### 2）  Singular can generate the ssh certificate key pair locally and automatically deploys the public key into virtual machine on a cloud. Then you can operate the virtual machine without a password. 
For example:  
  
```
$ singular configure --help
[Singular] (1/2) APIkey you have generated to access the public cloud API.
[Singular] (2/2) By using key-certificate pairs, you could access the Linux server without typing password.

$ singular configure
Cloud Access Key ID [None]: 6f2671de0d70ee5048379d16c0d0405df4a720ced263ffb35f67aded4834f321
Cloud SSHkey Path [./usr/singular/rsd_id.pub]: ./usr/singular/rsd_custom_name.pub
[Singular] API Server key is pass validation from cloud server.
[Singular] Generated Certificate Authority SSHkey and certificate.
[Singular] Created keys and certificates in "./usr/singular/rsd_custom_name.pub"
```
Note: API server key is required for authentication while calling API.
Note: Each step of the virtual machine operation depends on if your local private key matches virtual machine public key. It is more secure compared to the use of account password


### Getting Started

##### (1/3) Create Kubernetes cluster automatically with wizard. Configure cluster size and node setting with the singular, a YAML file will be generated.
```
$ singular deploy cluster 
    Master Count [1]: 3
    Node Count [2]: 3
    Memory Size [512]: 1024
    Region  [sfo]: 
    System Version [ubuntu-17-04-x64]: 
    NAME                   STATUS     AGE     REGION  
    ubuntu-master-1        Ready      1024M      sfo   
    ubuntu-master-2        Ready      1024M      sfo   
    ubuntu-master-3        Ready      1024M      sfo         
    ubuntu-node-1          Ready      1024M      sfo   
    ubuntu-node-2          Ready      1024M      sfo   
    ubuntu-node-3          Ready      1024M      sfo
    [singular]Are you sure you want to continue creating?[yes/no]
```
##### (2/3)  By calling call the public cloud API, singular can build your virtual machine nodes and retrieve the nodes information list.
```
    [singular]Are you sure you want to continue creating?[yes/no] yes
    NAME                   STATUS     PROGRESS          IP
    ubuntu-master-1        Ready      100%        138.68.14.197
    ubuntu-master-2        Ready      100%        138.68.14.198
    ubuntu-master-3        NoReady     80%              -
    [singular] Are you sure you want to continue deploying?[yes/no]
```
##### (3/3)  According the list, singular can download Kubernetes binary files to each node, and start deployment.
```
[singular]Are you sure you want to continue deploying?[yes/no] yes
NAME                 Donwload      Deploy       STATUS
ubuntu-master-1        100%         100%        SUCCEED
ubuntu-master-2        100%         100%        FAILED
ubuntu-master-3        80%          0%            -
```
Note: You could manually configure YAML file, and then execute deploy to setup and install. However, without the configuration file, part of information will be lost after singular destroyed, such as the path for API key and SSH certification.


### DESCRIPTION COMMAND & OPTION    
```
Available Commands:
configure   Configure your APIkey and API Server key and SSH certification of Kubernetes cluster with the wizard.
deploy      To start a new Kubernetes cluster deploying and running each service.
options:
    --config=~/usr/singular/config.yaml  setting custom singular path of config.
    --security                           Generate Kubernetes certificate
```
## Using singular with a configuration file
##### It’s possible to configure the Singular with a configuration file instead of command line flags, and some more advanced features may only be available as configuration file options. 

###Sample Configuration

```
cluster_config:
    SSHkey:     ""
    APIkey:     ""
    EtcdNet: "/kube/network"
    Security: "yes"

vm_config:
    Memory Size:     "1024mb"
    Region:     "sfo2"
    System Version:    "ubuntu-17-04-x64"
cluster_download:
    cloud_controller_manager : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/cloud-controller-manager"
    hyperkube : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/hyperkube"
    kube_aggregator : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kube-aggregator"
    kube_APIserver : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kube-APIserver"
    kube_controller_manager : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kube-controller-manager"
    kube_proxy : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kube-proxy"
    kube_scheduler : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kube-scheduler"
    kubeadm : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kubeadm"
    kubectl : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kubectl"
    kubefed : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kubefed"
    kubelet : "https://storage.googleAPIs.com/containerops-release/Kubernetes/1.6.2/kubelet"
```