# CLI

## Build

```shell
docker build -t hub.opshub.sh/binary/v1/containerops/component/binary/php/cli:0.1 .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/wp-cli/wp-cli.git file=./bin/wp composer=true" hub.opshub.sh/binary/v1/containerops/component/binary/php/cli:0.1
```

## Options

Required:

- git-url
- file
