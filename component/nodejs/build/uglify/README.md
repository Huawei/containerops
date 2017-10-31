# UGLIFY

## Build

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-uglify:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/mishoo/tweeg.js.git file=tweeg.js output=bundle.js" hub.opshub.sh/containerops/build-nodejs-uglify:latest
```

## Options

Required:

- git-url
- file
- output
