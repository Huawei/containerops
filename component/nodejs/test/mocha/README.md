## Node.js Code Test Component MOCHA

### What's the Component?

This image is node.js runtime image, used to testing your node.js project.

Mocha is a feature-rich JavaScript test framework running on Node.js and in the browser, making asynchronous testing simple and fun. Mocha tests run serially, allowing for flexible and accurate reporting, while mapping uncaught exceptions to the correct test cases. Hosted on GitHub.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/test-nodejs-mocha:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/expressjs/express.git \
    require=test/support/env \
    reporter=spec \
    bail=true \
    check-leaks=test/ \
    path=test/acceptance/" \
    hub.opshub.sh/containerops/test-nodejs-mocha:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `require` test requires
- `reporter` report type
- `bail` is bailed 
- `check-leaks` is checking leaks
- `path` test file path

### Versions 1.0.0