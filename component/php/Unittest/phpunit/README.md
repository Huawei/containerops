## PHP Code Unittest Component PHPUNIT

### What's the Component?

This image is php runtime image, used for testing your php project. 

PHPUnit is a programmer-oriented testing framework for PHP. It is an instance of the xUnit architecture for unit testing frameworks.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/unittest-php-phpunit:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/sebastianbergmann/phploc.git \
    --configuration=phpunit.xml \
    composer=true" \
    hub.opshub.sh/containerops/unittest-php-phpunit:latest
```

### Parameters 

Required:

- `git-url` where your code is located

Optional:

- `bootstrap` bootstrap=<file>          A "bootstrap" PHP file that is run before the tests.
- `composer` composer=true/false
- `include-path` include-path=<path(s)>    Prepend PHP's include_path with given path(s).
- `configuration` configuration=<file>   Read configuration from XML file.

### Versions 1.0.0