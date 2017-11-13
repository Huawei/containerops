## Node.js Code Document Component JSDOC

### What's the Component?

This image is node.js runtime image, used to generate document for your node.js project.

JSDoc 3 is an API documentation generator for JavaScript, similar to Javadoc or phpDocumentor. You add documentation comments directly to your source code, right alongside the code itself. The JSDoc tool will scan your source code and generate an HTML documentation website for you.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/document-nodejs-jsdoc:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/gitgrimbo/jsdoc3-examples.git \
    file=js/Book.js \
    config=conf.json" \
    hub.opshub.sh/containerops/document-nodejs-jsdoc:latest
```

## Options

Required:

- git-url
- file
- output

### Parameters 

Required:

- `git-url` where your code is located
- `file` targe files
- `output` output path

### Versions 1.0.0