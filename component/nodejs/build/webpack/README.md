## Node.js Code Build Component Webpack

### What's the Component?

This image is node.js runtime image, used for build your node.js project.

A bundler for javascript and friends. Packs many modules into a few bundled assets. Code Splitting allows to load parts for the application on demand. Through "loaders," modules can be CommonJs, AMD, ES6 modules, CSS, Images, JSON, Coffeescript, LESS, ... and your custom stuff.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/build-nodejs-webpack:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/yangkghjh/try_react.git \
    config=webpack.production.config.js" \
    hub.opshub.sh/containerops/build-nodejs-webpack:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `config` the webpack config file

### Versions 1.0.0