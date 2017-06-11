## Build, Test And Release CoreDNS

```bash
docker build -t docker.io/containerops/cncf-demo-coredns .
```


```bash
docker run --env CO_DATA="coredns=https://github.com/coredns/coredns.git action=build" docker.io/containerops/cncf-demo-coredns:latest
```

```dockerfile
FROM docker.io/containerops/golang:1.8.3
MAINTAINER Quanyi Ma <genedna@gmail.com>
USER root
RUN apt-get update && apt-get install -y gcc make g++
ENV PATH $PATH:$GOPATH/src/github.com/Huawei/containerops
RUN mkdir -p $GOPATH/src/github.com/Huawei/containerops
ADD codes/*.go $GOPATH/src/github.com/Huawei/containerops/
WORKDIR $GOPATH/src/github.com/Huawei/containerops
RUN go build coredns.go
WORKDIR $GOPATH/src/github.com/coredns/coredns
CMD coredns
```