# LESS

## Build

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-less:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/yanlibo2013/less.git file=less/mooc3.1.2/index.less output=build/index.css" hub.opshub.sh/containerops/build-nodejs-less:latest
```

## Options

Required:

- git-url
- file
- output
