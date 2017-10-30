# NPM

## Build

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-grunt:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/gruntjs/grunt-contrib-htmlmin.git action=test" hub.opshub.sh/containerops/build-nodejs-grunt:latest
```

## Options

Required:

- git-url
- action
