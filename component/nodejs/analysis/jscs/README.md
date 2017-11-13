## Node.js Code Analysis Component JSCS

### What's the Component?

This image is node.js runtime image, used for analysis your node.js coding style. 

JSCS is a code style linter and formatter for your style guide

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/analysis-nodejs-jscs:latest .
```

### Component Usage

```shell
docker run \
    --env  CO_DATA=" \
    git_url=https://github.com/spyl94/react-brunch-demo.git \
    path=brunch-config.js \
    preset=airbnb" 
    hub.opshub.sh/containerops/analysis-nodejs-jscs:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `path` code path
- `preset` code style

### Versions 1.0.0