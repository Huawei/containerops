# NPM

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-brunch:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/brunch/with-es6.git action=build" hub.opshub.sh/containerops/dependence-nodejs-brunch:latest
```

## Options

Required:

- git-url
- action
