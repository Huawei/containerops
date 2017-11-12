## Node.js Code Test Component JEST

### What's the Component?

This image is node.js runtime image, used to testing your node.js project.

Jest is used by Facebook to test all JavaScript code including React applications. One of Jest's philosophies is to provide an integrated "zero-configuration" experience. We observed that when engineers are provided with ready-to-use tools, they end up writing more tests, which in turn results in more stable and healthy code bases.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/test-nodejs-jest:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/ant-design/ant-design.git \
    config=.jest.js" \
    hub.opshub.sh/containerops/test-nodejs-jest:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `config` config file

### Versions 1.0.0