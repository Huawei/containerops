## PHP Code Compile Component PHAR

### What's the Component?

This image is php runtime image, used for compile your php project. 

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/compile-php-phar:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/SegmentFault/phar-sample.git \
    entry-file=build.php" \
    hub.opshub.sh/containerops/compile-php-phar:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `entry-file` entry file

Optional:

- `composer` true/false

### Versions 1.0.0