## Node.js Code Build Component Brunch

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

Brunch lets you focus on what matters most — solving real problems instead of messing around with the glue.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-brunch:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/brunch/with-es6.git \
    action=build" \
    hub.opshub.sh/containerops/build-nodejs-brunch:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `action` action you want to run

### Versions 1.0.0