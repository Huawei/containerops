# Singular

### Application


```
$ singular
Usage: singular [OPTIONS] COMMAND [arg...]
       singular [ --help | -v | --version ]

Singular, the Kubernetes deployment and operations tools.
```
#### To automatically deploy Kubernetes, you could simply follow below steps:
### Before You Start
Before using the singular, you need to tell it about your public cloud credentials.

```
$ singular config 
Welcome to singular!
Which one is your cloud compute provider: 
    [1] DigitalOcean 
    [2] Amazon Web Services
    [3] Google Cloud Platfrom
Input item number : 1

DigitalOcean is your provider.Please input your Cloud API Access Key. [Your could find it from your Account with following website.]
https://cloud.digitalocean.com/settings/api/tokens.

API Access Key ID : 6f267**********************321D34
API Server key is pass validation from cloud server.

Cloud SSHkey Path : /etc/singular/rsd_custom_name.pub
Generated Certificate Authority SSHkey and certificate.
Created SSHkey and certificate successfully. 
Setting up SSH keys on your cloud account.
Now,You can create new virtual machine  with an SSH key already set on them.
Congratulations, Let us get deployd up and running and build your container cluster !

```

### Getting Started

```

$ singular deploy k8s 
Let us start to deploy kubernetes cluster wish singular!
Nodes configuration：
Number of Master in cluster : 3
Number of Node in cluster : 3

Master[3]/Node[3]. Now, set virtual machine size :
[1] 2CPUs 4GB  Memory  60GB SSDdisk  4TB transfer
[2] 4CPUs 8GB  Memory  80GB SSDdisk  5TB transfer
[3] 8CPUs 16GB Memory 160GB SSDdisk  6TB transfer
Input item number : 1

The Nodes Size same as Mastes[Yes/no]:n
Node's virtual machine size :
[1] 2CPUs 4GB  Memory  60GB SSDdisk  4TB transfer
[2] 4CPUs 8GB  Memory  80GB SSDdisk  5TB transfer
[3] 8CPUs 16GB Memory 160GB SSDdisk  6TB transfer
Input item number : 2

Done. At last, select virtual machine region:
[1]New York [2]San Francisco [3]Singapore [4]Frankfurt 
Input item number : 2

Based on your selection, generate a list for you. 

---------------------------------------------------------
| HOST NAME                      | Price /Monthly &Hour |
---------------------------------------------------------
| k8s-master-ubuntu-4gb-NYC1-01  | $40/mo  $0.060 /hour |
| k8s-master-ubuntu-4gb-NYC1-02  | $40/mo  $0.060 /hour |
| k8s-master-ubuntu-4gb-NYC1-03  | $40/mo  $0.060 /hour |     
| k8s-node-ubuntu-8gb-NYC1-04    | $80/mo  $0.119 /hour |
| k8s-node-ubuntu-8gb-NYC1-05    | $80/mo  $0.119 /hour |
| k8s-node-ubuntu-8gb-NYC1-06    | $80/mo  $0.119 /hour |
---------------------------------------------------------
Add up $360/mo or $0.537/hour.
Are you sure you want to continue creating?[Yes/no]:
```
```
Are you sure you want to continue creating?[Yes/no]:y
This will download and install the official compiler of kubernetes for the Cluster.

---------------------------------------------------------------
| NAME           		|    STATUS   |      IP       | 
---------------------------------------------------------------
| k8s-master-ubuntu-4gb-NYC1-01 | Success     | 138.68.14.191 | 
| k8s-master-ubuntu-4gb-NYC1-02 | Installing  | 138.68.14.192 |  
| k8s-master-ubuntu-4gb-NYC1-03 | installing  | 138.68.14.193 |   
| k8s-node-ubuntu-8gb-NYC1-04   | Downloading | 138.68.14.197 |   
| k8s-node-ubuntu-8gb-NYC1-05   | VM Created  | 138.68.14.198 | 
| k8s-node-ubuntu-8gb-NYC1-06   | VM Creating |        -      | 
---------------------------------------------------------------
Kubernetes is installed now. Great!

```


### DESCRIPTION COMMAND & OPTION    
```

$ singular deploy cluster --master-number 2 --node-number 3 --mSize 4 --region sfo --storage 100

Available Commands:
config      Configure your APIkey and API Server key and SSH certification of Kubernetes cluster with the wizard.
deploy      To start a new Kubernetes cluster deploying and running each service.
options:
--config=~/etc/singular/config.yaml  setting custom singular path of config.
--security                           Generate Kubernetes certificate
--master-number Master nodes number of kubernetes cluster
--node-number   Slaves nodes number of kubernetes cluster
--mSize         Memory size of virtual machine
--region        Location region of virtual machine
--storage       Storage volume region of virtual machine
```
## Using singular with a configuration file

### Sample Configuration

```
cluster_config:
	SSHkey: ""
	APIkey: ""
	EtcdNet: "/kube/network"
	Security: "True"
	Private networking:"True"
        Region: "sfo2"

vm_config:
    Memory Size:"8"
    System Version:    "ubuntu-17-04-x64"
    Block storage:"100"
```
