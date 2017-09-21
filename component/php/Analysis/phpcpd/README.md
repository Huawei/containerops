# PHPCPD

## Build

```shell
docker build -t hub.opshub.sh/containerops/analysis-php-phpcpd:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git" hub.opshub.sh/containerops/analysis-php-phpcpd:latest
```

## Options

Required:

- git-url

Optional:

- path
- names
- names-exclude
- regexps-exclude
- exclude
- min-lines
- min-tokens

```shell
Options:
    path=.                           Files and directories to analyze
    names=NAMES                      A comma-separated list of file names to check [default: ["*.php"]]
    names-exclude=NAMES-EXCLUDE      A comma-separated list of file names to exclude
    regexps-exclude=REGEXPS-EXCLUDE  A comma-separated list of paths regexps to exclude (example: "#var/.*_tmp#")
    exclude=EXCLUDE                  Exclude a directory from code analysis (must be relative to source) (multiple values allowed)
    min-lines=MIN-LINES              Minimum number of identical lines [default: 5]
    min-tokens=MIN-TOKENS            Minimum number of identical tokens [default: 70]
```