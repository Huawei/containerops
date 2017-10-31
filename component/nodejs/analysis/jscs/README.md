# JSCS

## Build

```shell
docker build -t hub.opshub.sh/containerops/analysis-nodejs-jscs:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/spyl94/react-brunch-demo.git path=brunch-config.js preset=airbnb" hub.opshub.sh/containerops/analysis-nodejs-jscs:latest
```

## Options

Required:

- git-url
- path
- preset
