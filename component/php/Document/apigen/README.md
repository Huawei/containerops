## PHP Code Document Component APIGEN

### What's the Component?

This image is php runtime image, used for creating document for your php project. 

ApiGen is the simplest, the easiest to use and the most modern api doc generator. It is all PHP 7.1 features ready easy to extend with own Finder, Annotation Subscriber or even Generator.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/document-php-apigen:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/TIGERB/easy-php.git \
    path=app \
    destination=docs" \
    hub.opshub.sh/containerops/document-php-apigen:latest
```

### Parameters 

Required:

- `git-url` where your code is located

Optional:

- `path` 
- `exclude` List of directories to exclude, separated by a comma (,)
- `ignore-annotations`

### Versions 1.0.0