# MOCHA

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-mocha:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/expressjs/express.git require=test/support/env reporter=spec bail=true check-leaks=test/ path=test/acceptance/" hub.opshub.sh/containerops/dependence-nodejs-mocha:latest
```

## Options

Required:

- git-url
- config
