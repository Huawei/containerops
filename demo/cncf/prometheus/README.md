## Build, Test And Release Prometheus

```bash
docker build -t docker.io/containerops/cncf-demo-prometheus .
```


```bash
docker run --env CO_DATA="prometheus=https://github.com/prometheus/prometheus.git action=test" docker.io/containerops/cncf-demo-prometheus:latest
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
RUN go build prometheus.go
WORKDIR $GOPATH/src/github.com/prometheus/prometheus
CMD prometheus
```