# COMPOSER

## Build

```shell
docker build -t hub.opshub.sh/binary/v1/containerops/component/binary/php/component-composer: .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/sebastianbergmann/phploc.git entry-file=build.php" hub.opshub.sh/binary/v1/containerops/component/binary/php/component-composer:0.1
```

## Options

Required:

- git-url
