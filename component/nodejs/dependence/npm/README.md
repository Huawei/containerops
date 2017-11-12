## Node.js Code Dependence Component NPM

### What's the Component?

This image is node.js runtime image, used for manage your node.js project's dependence.

Npm is the package manager for JavaScript and the world’s largest software registry. Discover packages of reusable code — and assemble them in powerful new ways.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-npm:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/yangkghjh/try_react.git" \
    hub.opshub.sh/containerops/dependence-nodejs-npm:latest
```

### Parameters 

Required:

- `git-url` where your code is located

### Versions 1.0.0