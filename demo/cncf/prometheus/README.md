## Build, Test And Release Prometheus

```bash
docker build -t containerops/cncf-demo-prometheus .
```


```bash
docker run --env CO_DATA="prometheus=github://github.com/prometheus/prometheus.git action=test" containerops/cncf-demo-prometheus:latest
```
