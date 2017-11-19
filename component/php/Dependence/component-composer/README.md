## PHP Code Dependence Component Composer

### What's the Component?

This image is php runtime image, used for manage your php project's dependence. 

Dependency Manager for PHP.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/dependence-php-composer:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/sebastianbergmann/phploc.git" \
    hub.opshub.sh/containerops/dependence-php-composer:latest
```

### Parameters 

Required:

- `git-url` where your code is located

### Versions 1.0.0