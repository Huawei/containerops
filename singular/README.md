## Singular

Singular design for deploy and operation ContainerOps platform, mostly focus on Kubernetes, Prometheus, others in the Cloud Native technology stack. We are trying to build all stack cross cloud, OpenStack, bare metals.

Singular don't use any other deploy tools like _kubeadm_, deploy everything in a hard way instead. Singular providers templates for service so could deploy any version. Singular deploys development versions of **Kubernetes**, **CoreDNS**, others in **CNCF** CI demo.

### Deployment Template

Singular uses a **YAML** file as deploy template, it describes the architecture of the cluster. It don't only deploy ContainerOps modules, also use to deploy Kubernetes, others in Cloud Native stack and some common software.

#### SSH Key

When deploy infrastructures, **Singular** need to _SSH_ to virtual machine or bare metal.

1. At least provide SSH private key, **Singular** will create public key from it. Generate ssh key follow Github document - [How to generate SSH Key](https://help.github.com/articles/generating-a-new-ssh-key-and-adding-it-to-the-ssh-agent).
2. If no SSH key files in deployment template, **Singular** will create _SSH_ pair key files in default folder(**_$HOME/.containerops/ssh_**) and name(**_id_rsa.pub_** and **_id_rsa_**).

##### Deploy Command

```
singular deploy template /tmp/deploy.yml  --verbose --timestamp
```