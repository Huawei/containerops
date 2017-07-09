## php Runtime Docker Image

### What's the image?

This image is php runtime image, used for build php application. 

### How to build the image?

Use the `docker build` command build the image, and `php_version` is you build 

```
docker build -t containerops/php:7.1.4 --build-arg php_version=7.1.4  .
```

### Versions


