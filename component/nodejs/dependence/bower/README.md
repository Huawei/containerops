# NPM

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-bower:latest .
```

## Run

```shell
docker run --env CO_DATA="git_url=https://github.com/WildDogTeam/demo-js-wildchat.git" hub.opshub.sh/containerops/dependence-nodejs-bower:latest
```

## Options

Required:

- git-url
