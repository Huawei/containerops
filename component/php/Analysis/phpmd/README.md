## PHP Code Analysis Component PHPMD

### What's the Component?

This image is php runtime image, used for analysis your php coding style. 

PHPMD is a spin-off project of PHP Depend and aims to be a PHP equivalent of the well known Java tool PMD. PHPMD can be seen as an user friendly frontend application for the raw metrics stream measured by PHP Depend. http://phpmd.org

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/analysis-php-phpmd:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/TIGERB/easy-php.git" \
    hub.opshub.sh/containerops/analysis-php-phpmd:latest
```

### Parameters 

Required:

- `git-url` where your code is located

Optional:

- `path` A php source code filename or directory. Can be a comma-separated string
- `formats` A report format.Available formats: xml, text, html.
- `ruleset` A ruleset filename or a comma-separated string of rulesetfilenames.Available rulesets: cleancode, codesize, controversial, design, naming, unusedcode.
- `minimumpriority` =true/false.rule priority threshold; rules with lower priority than this will not be used
- `exclude` comma-separated string of patterns that are used to ignore directories
- `suffixes` comma-separated string of valid source code filename extensions, e.g. php,phtml
- `strict` also report those nodes with a @SuppressWarnings annotation
- `ignore-violations-on-exit` will exit with a zero code, even if any violations are found

### Versions 1.0.0