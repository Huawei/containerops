## Build, Test And Publish Kubernetes With Bazel

```bash
docker build -t containerops/cncf-demo-kubernetes .
```


```bash
docker run --env CO_DATA=`kubernetes=https://github.com/kubernetes/kubernetes.git action=build` containerops/cncf-demo-kubernetes:latest
```
