# JSHINT

## Build

```shell
docker build -t hub.opshub.sh/containerops/analysis-nodejs-jshint:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/jshint/jshint.git path=src" hub.opshub.sh/containerops/analysis-nodejs-jshint:latest
```

## Options

Required:

- git-url
- config
