## Node.js Code Build Component Browserify

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

Browserify lets you require('modules') in the browser by bundling up all of your dependencies.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-browserify:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/mattdesl/browserify-example.git \
    output=bundle.js \
    file=index.js" \
    hub.opshub.sh/containerops/build-nodejs-browserify:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `output` output path
- `file` file path

### Versions 1.0.0