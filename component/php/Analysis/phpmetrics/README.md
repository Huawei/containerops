## PHP Code Analysis Component PHPMETRICS

### What's the Component?

This image is php runtime image, used for analysis your php coding style. 

PhpMetrics provides metrics about PHP project and classes, with beautiful and readable HTML report.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/analysis-php-phpmetrics:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/TIGERB/easy-php.git" \
    hub.opshub.sh/containerops/analysis-php-phpmetrics:latest
```

### Parameters 

Required:

- `git-url` where your code is located

Optional:

- `path`
- `exclude` List of directories to exclude, separated by a comma (,)
- `ignore-annotations`

### Versions 1.0.0