## nodejs Runtime Docker Image

### What's the image?

This image is nodejs runtime image, used for build nodejs application. 

### How to build the image?

Use the `docker build` command build the image, and `node_version` is you build 

```
docker build -t containerops/node:8.1.3 --build-arg node_version=8.1.3 .
```

### Versions


