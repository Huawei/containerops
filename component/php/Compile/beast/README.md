## PHP Code Compile Component BEAST

### What's the Component?

This image is php runtime image, used for compile your php project. 

PHP source code encrypt module.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/compile-php-beast:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/TIGERB/easy-php.git" \
    hub.opshub.sh/containerops/compile-php-beast:latest
```

### Parameters 

Required:

- `git-url` where your code is located

Optional:

- `composer` true/false

### Versions 1.0.0