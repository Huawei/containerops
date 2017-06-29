## Golang Runtime Docker Image

### What's the image?

This image is Golang runtime image, used for build Golang application. 

### How to build the image?

Use the `docker build` command build the image, and `go_version` is you build 

```
docker build -t containerops/golang:1.8.1 --build-arg go_version=1.8.1  .
```

### Versions


