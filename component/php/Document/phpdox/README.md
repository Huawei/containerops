## PHP Code Document Component PHPDOX

### What's the Component?

This image is php runtime image, used for creating document for your php project. 

phpDox is a documentation generator for PHP projects. This includes, but is not limited to, API documentation. The main focus is on enriching the generated documentation with additional details like code coverage, complexity information and more.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/document-php-phpdox:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/theseer/phpdox.git" \
    hub.opshub.sh/containerops/document-php-phpdox:latest
```

### Parameters 

Required:

- `git-url` where your code is located

### Versions 1.0.0