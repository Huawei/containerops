# JEST

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-jest:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/ant-design/ant-design.git config=.jest.js" hub.opshub.sh/containerops/dependence-nodejs-jest:latest
```

## Options

Required:

- git-url
- config
