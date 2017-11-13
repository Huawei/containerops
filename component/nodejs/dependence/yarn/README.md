## Node.js Code Dependence Component Yarn

### What's the Component?

This image is node.js runtime image, used for manage your node.js project's dependence.

FAST, RELIABLE, AND SECURE DEPENDENCY MANAGEMENT.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-yarn:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/yangkghjh/try_react.git" \
    hub.opshub.sh/containerops/dependence-nodejs-yarn:latest
```

### Parameters 

Required:

- `git-url` where your code is located

### Versions 1.0.0
