## Node.js Code Analysis Component Eslint

### What's the Component?

This image is node.js runtime image, used for analysis your node.js coding style. 

ESLint is an open source project originally created by Nicholas C. Zakas in June 2013. Its goal is to provide a pluggable linting utility for JavaScript.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/analysis-nodejs-eslint:latest .
```

### Component Usage

```shell
docker run \
   --env CO_DATA=" \
   git_url=https://github.com/spyl94/react-brunch-demo.git \
   path=src" 
   hub.opshub.sh/containerops/analysis-nodejs-eslint:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `path` code path

### Versions 1.0.0