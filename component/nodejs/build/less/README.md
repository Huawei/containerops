## Node.js Code Build Component LESS

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

Less is a CSS pre-processor, meaning that it extends the CSS language, adding features that allow variables, mixins, functions and many other techniques that allow you to make CSS that is more maintainable, themeable and extendable.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-less:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/yanlibo2013/less.git \
    file=less/mooc3.1.2/index.less \
    output=build/index.css" \
    hub.opshub.sh/containerops/build-nodejs-less:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `file` file path
- `output` output path

### Versions 1.0.0