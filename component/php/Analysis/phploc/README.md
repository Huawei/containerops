# PHPLOC

## Build

```shell
docker build -t hub.opshub.sh/binary/v1/containerops/component/binary/php/phploc:0.1 .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git exclude=public" hub.opshub.sh/binary/v1/containerops/component/binary/php/phploc:0.1
```

## Options

Required:

- git-url

Optional:

- path
- exclude
- names
- names-exclude

```shell
names=NAMES                  A comma-separated list of file names to check [default: ["*.php"]]
names-exclude=NAMES-EXCLUDE  A comma-separated list of file names to exclude
exclude=EXCLUDE              Exclude a directory from code analysis (multiple values allowed)
```