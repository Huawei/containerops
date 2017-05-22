# Singular

### Application


```
$ singular
Usage: singular [OPTIONS] COMMAND [arg...]
       singular [ --help | -v | --version ]

The kubernetes deployment and operations tools.
To automatically deply kubernetes, you could simply follow below steps:

```
### Getting Started

##### 1）  Configure cluster size and node with singular, a yaml file will be generated.
```
$ singular config master --count 2 --mSize 512 --region sfo --slug buntu-17-04-x64
NAME                   STATUS     MSize     REGION 
ubuntu-master-1        Ready      512M        sfo   
ubuntu-master-2        Ready      512M        sfo   
$ singular config node   --count 3 --mSize 1024 --region sfo --slug buntu-17-04-x64
NAME                   STATUS     AGE     REGION       
ubuntu-minion-1        Ready      1024M      sfo   
ubuntu-minion-2        Ready      1024M      sfo   
ubuntu-minion-3        Ready      1024M      sfo   

```
##### 2）  By calling call the public cloud API, singular can build your vm node and retrieve the node information list.
```
$ singular deploy master 
NAME                   STATUS     PROGRESS          IP
ubuntu-master-1        Ready      100%        138.68.14.197
ubuntu-master-2        Ready      100%        138.68.14.198
ubuntu-master-3        NoReady     80%              -
```
##### 3）  According the list, singular can download kubernetes binary files to each node, and start deployment with the yaml file generated in step 1.
```
$ singular deploy master 
NAME                 Donwload      Deploy       STATUS
ubuntu-master-1        100%         100%        SUCCEED
ubuntu-master-2        100%         100%        FAILED
ubuntu-master-3        80%          0%            -
```

Note:
You could manually configure yaml file, and then execute setup to deploy and install. However, without the configuration file, part of information will be lost after singular destroyed, such as the path for api key and cert.
### Precondition
Before using a singular, you need to tell it about your AWS credentials. You can do this in several steps:

 
Note: api server key is required for authentication while call api.
For example:
 
##### 1）  Register an account for public cloud, and retrieve api server key and put it into the yaml file.
Note: api server key is required for authentication while call api.
For example:  
  
```
$ singular apikey  6f2671de0d70ee5048379d16c0d0405df4a720ced263ffb35f67aded4834f321
[singular] API Server key is ready.
```

##### 2）  Singular can generates the ssh certificate key pair locally and automatically deploys the public key into vm. Then you can operate the virtual machine without a password. 
Note:Each step of the virtual machine operation depends on if your local private key matches virtual machine public key. It is more secure compared to the use of account password
    
```
$ singular cerpath  ./usr/singular/
$ singular cerkey 
[singular] Generated Certificate Authority key and certificate.
[singular] Created keys and certificates in "/usr/singular/"
```
### DESCRIPTION COMMAND & OPTION
    
```
Available Commands:
config  Configure your nodes of kubernetes cluster
deploy  To start a new kubernetes cluster deploying and running each services
cluster  Get kubernetes cluster information and status
apikey	APIkey you have generated to access the public cloud API.
cerkey	Generated key-certificate pairs could help to access to the linux server without the need to type password.
		by using Generated key-certificate pairs, you could access to the linux server without typing password.
options:
		
		--config=~/etc/.singular/config.yaml setting custom singular path of config.
        --cerpath	Without CApath option ,the default value is /etc/.singular/id_rsa.pub
					Or you could type your custom path for generate file id_rsa and id_rsa.pub
		--master|slave                       Create master or slave nodes
		--security                           Generate kubernetes certificate
		--privtenet       					 Privte network for your cluster
		--count =3   <value>		    	 Number of nodes in cluster
        --mSize =512	 <1024|2048|>        Node memory Size
        --region =sfo    sfo|nyc			 Cluster's localization of the region
        --slug=ubuntu-17-04-x64  <value>     System version
        --pull             Download Kubernetes binaries without install.               
        
```
## Using singular with a configuration file
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