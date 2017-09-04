# APIGEN

## Build

```shell
docker build -t hub.opshub.sh/binary/v1/containerops/component/binary/php/apigen:0.1 .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git path=app destination=docs" hub.opshub.sh/binary/v1/containerops/component/binary/php/apigen:0.1
```

## Options

Required:

- git-url

Optional:

- path
- exclude
- ignore-annotations

```shell
exclude=<directory>               List of directories to exclude, separated by a comma (,)
```