# CLI

## Build

```shell
docker build -t hub.opshub.sh/containerops/base-php-cli:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/wp-cli/wp-cli.git file=./bin/wp composer=true" hub.opshub.sh/containerops/base-php-cli:latest
```

## Options

Required:

- git-url
- file
