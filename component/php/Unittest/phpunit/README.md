# PHPUNIT

## Build

```shell
docker build -t hub.opshub.sh/binary/v1/containerops/component/binary/php/phpunit:0.1 .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/sebastianbergmann/phploc.git --configuration=phpunit.xml composer=true" hub.opshub.sh/binary/v1/containerops/component/binary/php/phpunit:0.1
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