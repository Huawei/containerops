# NPM

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-npm:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/yangkghjh/try_react.git" hub.opshub.sh/containerops/dependence-nodejs-npm:latest
```

## Options

Required:

- git-url
