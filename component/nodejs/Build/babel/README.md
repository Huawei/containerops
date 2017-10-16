# NPM

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-babel:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/ronen-e/babel-demo.git file=src/app.js o=dist/app.js" hub.opshub.sh/containerops/dependence-nodejs-babel:latest
```

## Options

Required:

- git-url
- file
- o
