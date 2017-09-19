# NPM

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-gulp:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/hjzheng/gulp-AngularJS1.x-seed.git action=build" hub.opshub.sh/containerops/dependence-nodejs-gulp:latest
```

## Options

Required:

- git-url
- action
