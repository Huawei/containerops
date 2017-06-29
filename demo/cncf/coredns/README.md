## Build, Test And Release CoreDNS

### How to Build Component
```bash
docker build -t docker.io/containerops/cncf-demo-coredns .
```

### How to Run The Component
```bash
docker run --env CO_DATA="coredns=https://github.com/coredns/coredns.git action=build release=hub.opshub.sh/containerops/cncf-demo/demo" docker.io/containerops/cncf-demo-coredns:latest
```

### Dockerfile
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