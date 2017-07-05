## Build, Test And Release Kubernetes With Bazel

```bash
docker build -t hub.opshub.sh/containerops/cncf-demo-kubernetes .
```


```bash
docker run --env CO_DATA="kubernetes=https://github.com/kubernetes/kubernetes.git action=build release=hub.opshub.sh/containerops/cncf-demo/demo" hub.opshub.sh/containerops/cncf-demo-kubernetes:latest
```

```dockerfile
FROM hub.opshub.sh/containerops/bazel:latest
MAINTAINER Quanyi Ma <genedna@gmail.com>

USER root
RUN apt-get update && apt-get install -y git make python python-dev python-pip python-virtualenv
RUN curl -sSL https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz -o /tmp/go.tar.gz && \
  echo "1862f4c3d3907e59b04a757cfda0ea7aa9ef39274af99a784f5be843c80c6772  /tmp/go.tar.gz" | sha256sum -c - && \
  tar -C /var/opt -xzf /tmp/go.tar.gz && \
  rm /tmp/go.tar.gz && \
  mkdir -p /var/opt/gopath && \
  chmod -R 777 /var/opt/gopath
ENV GOROOT /var/opt/go
ENV GOPATH /var/opt/gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH:/bin:$GOPATH/src/github.com/Huawei/containerops
RUN mkdir -p $GOPATH/src/github.com/Huawei/containerops
ADD codes/*.go $GOPATH/src/github.com/Huawei/containerops/
WORKDIR $GOPATH/src/github.com/Huawei/containerops
RUN go build k8s.go
WORKDIR $GOPATH/src/github.com/kubernetes/kubernetes
CMD k8s
```
