# COMPOSER

## Build

```shell
docker build -t hub.opshub.sh/containerops/dependence-php-composer:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/sebastianbergmann/phploc.git entry-file=build.php" hub.opshub.sh/containerops/dependence-php-composer:latest
```

## Options

Required:

- git-url
