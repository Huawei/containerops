## Node.js Code Analysis Component JSHINT

### What's the Component?

This image is node.js runtime image, used for analysis your node.js coding style. 

JSHint is a community-driven tool that detects errors and potential problems in JavaScript code. Since JSHint is so flexible, you can easily adjust it in the environment you expect your code to execute. JSHint is open source and will always stay this way.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/analysis-nodejs-jshint:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/jshint/jshint.git \
    path=src" \
    hub.opshub.sh/containerops/analysis-nodejs-jshint:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `path` code path

### Versions 1.0.0