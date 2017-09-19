# COFFEE

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-coffee:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/jonpliske/coffeescript_examples.git files=001_function_example.coffee output=bundle.js" hub.opshub.sh/containerops/dependence-nodejs-coffee:latest
```

## Options

Required:

- git-url
- file
- output
