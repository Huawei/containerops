# NPM

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-browserify:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/mattdesl/browserify-example.git output=bundle.js file=index.js" hub.opshub.sh/containerops/dependence-nodejs-browserify:latest
```

## Options

Required:

- git-url
- output
- file
