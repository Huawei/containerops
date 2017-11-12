## Node.js Code Build Component COFFEE

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

CoffeeScript is a little language that compiles into JavaScript. Underneath that awkward Java-esque patina, JavaScript has always had a gorgeous heart. CoffeeScript is an attempt to expose the good parts of JavaScript in a simple way.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-coffee:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/jonpliske/coffeescript_examples.git \
    files=001_function_example.coffee \
    output=bundle.js" \
    hub.opshub.sh/containerops/build-nodejs-coffee:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `file` file path
- `output` output path

### Versions 1.0.0