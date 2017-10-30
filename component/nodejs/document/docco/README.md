# DOCCO

## Build

```shell
docker build -t hub.opshub.sh/containerops/document-nodejs-docco:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/jonpliske/coffeescript_examples.git file=*.coffee" hub.opshub.sh/containerops/document-nodejs-docco:latest
```

## Options

Required:

- git-url
- file
