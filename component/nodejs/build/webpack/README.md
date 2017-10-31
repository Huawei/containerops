# NPM

## Build

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-webpack:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/yangkghjh/try_react.git config=webpack.production.config.js" hub.opshub.sh/containerops/build-nodejs-webpack:latest
```

## Options

Required:

- git-url
- config