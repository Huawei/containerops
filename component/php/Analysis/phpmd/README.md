# PHPMD

## Build

```shell
docker build -t hub.opshub.sh/containerops/analysis-php-phpmd:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/TIGERB/easy-php.git" hub.opshub.sh/containerops/analysis-php-phpmd:latest
```

## Options

Required:

- git-url

Optional:

- path
- formats
- ruleset
- minimumpriority
- exclude
- suffixes
- strict
- ignore-violations-on-exit

```shell
path: A php source code filename or directory. Can be a comma-separated string
formats: A report format.Available formats: xml, text, html.
ruleset: A ruleset filename or a comma-separated string of rulesetfilenames.Available rulesets: cleancode, codesize, controversial, design, naming, unusedcode.
minimumpriority=true/false: rule priority threshold; rules with lower priority than this will not be used
suffixes: comma-separated string of valid source code filename extensions, e.g. php,phtml
exclude: comma-separated string of patterns that are used to ignore directories
strict: also report those nodes with a @SuppressWarnings annotation
ignore-violations-on-exit: will exit with a zero code, even if any violations are found
```