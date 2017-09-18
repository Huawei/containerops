# PHPDOX

## Build

```shell
docker build -t hub.opshub.sh/containerops/document-php-phpdox:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/theseer/phpdox.git" hub.opshub.sh/containerops/document-php-phpdox:latest
```

## Options

Required:

- git-url
