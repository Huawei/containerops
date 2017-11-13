## Node.js Code Build Component Grunt

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

The JavaScript Task Runner.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-grunt:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/gruntjs/grunt-contrib-htmlmin.git \
    action=test" \
    hub.opshub.sh/containerops/build-nodejs-grunt:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `action` the action you want to run

### Versions 1.0.0