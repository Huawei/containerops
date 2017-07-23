## Deploy a kubernetes cluster with ContainerOps singular

```bash
docker build -t hub.opshub.sh/containerops/cncf-demo-singular .
```


```bash
docker run --env CO_DATA="\
	action=release \
	token=435a054fb1f81e439a63a608eddb67208a66cb11d6b7abeaa3d89aac777d4d1d \
	kube_apiserver_url=https://hub.opshub.sh/binary/v1/containerops/singular/binary/kube-apiserver/1.6.7 \
	kube_controllermanager_url=https://hub.opshub.sh/binary/v1/containerops/singular/binary/kube-controller-manager/1.6.7 \
	kube_scheduler_url=https://hub.opshub.sh/binary/v1/containerops/singular/binary/kube-scheduler/1.6.7 \
	kubectl_url=https://hub.opshub.sh/binary/v1/containerops/singular/binary/kubectl/1.6.7 \
	kubelete_url=https://hub.opshub.sh/binary/v1/containerops/singular/binary/kubelet/1.6.7 \
	kube_proxy_url=https://hub.opshub.sh/binary/v1/containerops/singular/binary/kube-proxy/1.6.7" \
hub.opshub.sh/containerops/cncf-demo-singular:latest

```

```dockerfile
FROM hub.opshub.sh/containerops/golang:1.8.3
MAINTAINER Zhen Ju <juzhenatpku@gmail.com>

USER root
RUN apt-get update && apt-get install -y gcc make g++
ENV PATH $PATH:$GOPATH/src/github.com/Huawei/containerops

RUN mkdir -p $GOPATH/src/github.com/Huawei/

WORKDIR $GOPATH/src/github.com/Huawei/
RUN git clone https://github.com/Huawei/containerops

# The essential dependencies for singular
RUN go get "github.com/cloudflare/cfssl/cli"
RUN go get "github.com/cloudflare/cfssl/cli/genkey"
RUN go get "github.com/cloudflare/cfssl/cli/sign"
RUN go get "github.com/cloudflare/cfssl/csr"
RUN go get "github.com/cloudflare/cfssl/initca"
RUN go get "github.com/cloudflare/cfssl/signer"
RUN go get "github.com/digitalocean/godo"
RUN go get "github.com/fernet/fernet-go"
RUN go get "github.com/logrusorgru/aurora"
RUN go get "github.com/mitchellh/go-homedir"
RUN go get "github.com/pkg/sftp"
RUN go get "github.com/spf13/cobra"
RUN go get "github.com/spf13/viper"
RUN go get "golang.org/x/crypto/ssh"
RUN go get "golang.org/x/net/context"
RUN go get "golang.org/x/oauth2"
RUN go get "gopkg.in/yaml.v2"

WORKDIR $GOPATH/src/github.com/Huawei/containerops/singular/
RUN go build

WORKDIR $GOPATH/src/github.com/Huawei/containerops/demo/cncf/singular/codes/
RUN go build -o ../singular-demo

WORKDIR $GOPATH/src/github.com/Huawei/containerops/demo/cncf/singular/
RUN cp $GOPATH/src/github.com/Huawei/containerops/singular/singular ./

CMD ./singular-demo
```
