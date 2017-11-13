## Node.js Code Build Component Gulp

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

gulp is a toolkit for automating painful or time-consuming tasks in your development workflow, so you can stop messing around and build something.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-gulp:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/hjzheng/gulp-AngularJS1.x-seed.git \
    action=build" \
    hub.opshub.sh/containerops/build-nodejs-gulp:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `action` the action you want to run

### Versions 1.0.0