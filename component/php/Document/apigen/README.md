# APIGEN

## Build

```shell
docker build -t hub.opshub.sh/containerops/document-php-apigen:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git path=app destination=docs" hub.opshub.sh/containerops/document-php-apigen:latest
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