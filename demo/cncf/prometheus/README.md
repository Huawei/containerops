## Build, Test And Release Prometheus

```bash
docker build -t docker.io/containerops/cncf-demo-prometheus .
```


```bash
docker run --env CO_DATA="prometheus=https://github.com/prometheus/prometheus.git action=test" docker.io/containerops/cncf-demo-prometheus:latest
```
