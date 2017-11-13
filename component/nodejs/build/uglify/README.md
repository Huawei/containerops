## Node.js Code Build Component UGLIFY

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

UglifyJS is a JavaScript compressor/minifier written in JavaScript.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-uglify:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA="git_url=https://github.com/mishoo/tweeg.js.git \
    file=tweeg.js \
    output=bundle.js" \
    hub.opshub.sh/containerops/build-nodejs-uglify:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `file` file path
- `output` output path

### Versions 1.0.0