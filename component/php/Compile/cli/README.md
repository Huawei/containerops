## PHP Code Compile Component CLI

### What's the Component?

This image is php runtime image, used for run php script. 

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/base-php-cli:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/wp-cli/wp-cli.git \
    file=./bin/wp \
    composer=true" \
    hub.opshub.sh/containerops/base-php-cli:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `file` the script

### Versions 1.0.0
