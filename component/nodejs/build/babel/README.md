## Node.js Code Build Component Babel

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

Babel is a JavaScript compiler.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-babel:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/ronen-e/babel-demo.git \
    file=src/app.js \
    o=dist/app.js" \
    hub.opshub.sh/containerops/build-nodejs-babel:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `file` file path
- `o` output path

### Versions 1.0.0