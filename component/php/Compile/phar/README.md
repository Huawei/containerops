# PHAR

## Build

```shell
docker build -t hub.opshub.sh/containerops/compile-php-phar:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/SegmentFault/phar-sample.git entry-file=build.php" hub.opshub.sh/containerops/compile-php-phar:latest
```

## Options

Required:

- git-url
- entry-file

Optional:

- composer

```shell
composer=true/false
```