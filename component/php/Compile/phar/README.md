# PHAR

## Build

```shell
docker build -t hub.opshub.sh/binary/v1/containerops/component/binary/php/phar:0.1 .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/SegmentFault/phar-sample.git entry-file=build.php" hub.opshub.sh/binary/v1/containerops/component/binary/php/phar:0.1
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