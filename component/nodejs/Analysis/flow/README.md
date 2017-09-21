# FLOW

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-flow:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/facebook/flow.git" hub.opshub.sh/containerops/dependence-nodejs-flow:latest
```

## Options

Required:

- git-url
