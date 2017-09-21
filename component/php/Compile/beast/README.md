# BEAST

## Build

```shell
docker build -t hub.opshub.sh/containerops/compile-php-beast:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git" hub.opshub.sh/containerops/compile-php-beast:latest
```

## Options

Required:

- git-url

Optional:

- composer

```shell
composer=true/false
```