## Node.js Code Dependence Component Bower

### What's the Component?

This image is node.js runtime image, used for manage your node.js project's dependence.

Bower is optimized for the front-end. If multiple packages depend on a package - jQuery for example - Bower will download jQuery just once. This is known as a flat dependency graph and it helps reduce page load.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-bower:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/WildDogTeam/demo-js-wildchat.git" \
    hub.opshub.sh/containerops/dependence-nodejs-bower:latest
```

### Parameters 

Required:

- `git-url` where your code is located

### Versions 1.0.0
