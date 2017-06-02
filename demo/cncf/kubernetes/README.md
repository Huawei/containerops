## Build, Test And Release Kubernetes With Bazel

```bash
docker build -t docker.io/containerops/cncf-demo-kubernetes .
```


```bash
docker run --env CO_DATA="kubernetes=https://github.com/kubernetes/kubernetes.git action=build" docker.io/containerops/cncf-demo-kubernetes:latest
```
