# PHPDOX

## Build

```shell
docker build -t hub.opshub.sh/binary/v1/containerops/component/binary/php/phpdox:0.1 .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/theseer/phpdox.git" hub.opshub.sh/binary/v1/containerops/component/binary/php/phpdox:0.1
```

## Options

Required:

- git-url
