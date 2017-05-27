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
Before using the singular, you need to tell it about your public cloud credentials. Two steps as following example:
For example:  
  
```
$ singular config 
Let's start 
Which one is your cloud compute provider: 
    [1] DigitalOcean 
    [2] Amazon Web Services
    [3] Google Cloud Platfrom
Input item number : 1

Input your Cloud API Access Key of DigitalOcean‘s Account.

API Access Key ID : 6f267**********************321D34
API Server key is pass validation from cloud server.

Singular can generate the ssh certificate key pair locally and automatically deploys the public key into virtual machine on a cloud. Then you can operate the virtual machine without a password.

Cloud SSHkey Path : /etc/singular/rsd_custom_name.pub

Generated Certificate Authority SSHkey and certificate.
Created SSHkey and certificate successfully. 
Setting up SSH keys on your cloud account.
Now,You can create new virtual machine  with an SSH key already set on them.
Congratulations, Let us get deployd up and running and build your container cluster !

```
Note: API server key is required for authentication while calling API.
Note: Each step of the virtual machine operation depends on if your local private key matches virtual machine public key. It is more secure compared to the use of account password

### Getting Started

##### (1/3) Create Kubernetes cluster automatically with wizard. Configure cluster size and node setting with the singular, a YAML file will be generated.
```

$ singular deploy k8s 
Master Count [1]: 3
Node Count [2]: 3
Virtual Machine Size :
[1] 4GB  2CPUs  60GB SSDdisk  4TB transfer
[2] 8GB  4CPUs  80GB SSDdisk  5TB transfer
[3] 16GB 8CPUs 160GB SSDdisk  6TB transfer

Select virtual machine region:
[1]New York [2]San Francisco [3]Amsterdam [4]Singapore[5]London
Input item number : 1

 ------------------------------------------------------
| NAME                             | SIZE  |  REGION |
-------------------------------------------------------
| k8s-master-ubuntu-4gb-nyc3-01    |  4G   |   sfo   |
| k8s-master-ubuntu-4gb-nyc3-01    |  4G   |   sfo   |
| k8s-master-ubuntu-4gb-nyc3-01    |  4G   |   sfo   |     
| k8s-node-ubuntu-4gb-nyc3-01      |  4G   |   sfo   |
| k8s-node-ubuntu-4gb-nyc3-01      |  4G   |   sfo   |
| k8s-node-ubuntu-4gb-nyc3-01      |  4G   |   sfo   |
-------------------------------------------------------

Are you sure you want to continue creating?[yes/no]
```
##### (2/3)  By calling call the public cloud API, singular can build your virtual machine nodes and retrieve the nodes information list.
```
Are you sure you want to continue creating?[Yes/no] y
-----------------------------------------------------------------
|       NAME       |       STATUS       |         IP            |
-----------------------------------------------------------------
| ubuntu-master-1  |       Ready        |      138.68.14.197    |
| ubuntu-master-2  |       Ready        |      138.68.14.198    |
| ubuntu-master-3  |       NoReady      |          -            |
-----------------------------------------------------------------

    Are you sure you want to continue deploying?[Yes/no]
```
##### (3/3)  According the list, singular can download Kubernetes binary files to each node, and start deployment.
```
Are you sure you want to continue deploying?[Yes/no] y
-----------------------------------------------------------------
| NAME              |    Donwload   |    Deploy    |    STATUS  | 
-----------------------------------------------------------------
| ubuntu-master-1   |      100%     |     100%     |    SUCCEED | 
| ubuntu-master-2   |      100%     |     100%     |    FAILED  | 
| ubuntu-master-3   |      80%      |      0%      |       -    | 
-----------------------------------------------------------------
```
Note: You could manually configure YAML file, and then execute deploy to setup and install. However, without the configuration file, part of information will be lost after singular destroyed, such as the path for API key and SSH certification.


### DESCRIPTION COMMAND & OPTION    
```
$ singular deploy cluster --master-count 2 --node-count 3 --mSize 512 --region sfo --sysversion ubuntu-17-04-x64

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

```