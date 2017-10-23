# PHPUNIT

## Build

```shell
docker build -t hub.opshub.sh/containerops/unittest-php-phpunit:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/sebastianbergmann/phploc.git --configuration=phpunit.xml composer=true" hub.opshub.sh/containerops/unittest-php-phpunit:latest
```

## Options

Required:

- git-url

Optional:

- bootstrap
- composer
- include-path
- configuration

```shell
bootstrap=<file>          A "bootstrap" PHP file that is run before the tests.
composer=true/false
include-path=<path(s)>    Prepend PHP's include_path with given path(s).
configuration=<file>   Read configuration from XML file.
```