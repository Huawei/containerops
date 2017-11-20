## PHP Code Analysis Component PHPLOC

### What's the Component?

This image is php runtime image, used for analysis your php coding style. 

`phploc` is a tool for quickly measuring the size and analyzing the structure of a PHP project.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/analysis-php-phploc:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/TIGERB/easy-php.git \
    exclude=public" \
    hub.opshub.sh/containerops/analysis-php-phploc:latest
```

### Parameters 

Required:

- `git-url` where your code is located

Optional:

- `path`
- `exclude` Exclude a directory from code analysis (multiple values allowed)
- `names`  A comma-separated list of file names to check [default: ["*.php"]]
- `names-exclude`  A comma-separated list of file names to exclude

### Versions 1.0.0