## Build, Test And Release CoreDNS

```bash
docker build -t docker.io/containerops/cncf-demo-coredns .
```


```bash
docker run --env CO_DATA="coredns=https://github.com/coredns/coredns.git action=build" docker.io/containerops/cncf-demo-coredns:latest
```
