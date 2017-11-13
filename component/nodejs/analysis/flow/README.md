## Node.js Code Analysis Component Flow

### What's the Component?

This image is node.js runtime image, used for analysis your node.js coding style. 

#### CODE FASTER.

Tired of having to run your code to find bugs? Flow identifies problems as you code. Stop wasting your time guessing and checking.

#### CODE SMARTER.

It's hard to build smart tools for dynamic languages like JavaScript. Flow understands your code and makes its knowledge available, enabling other smart tools to be built on top of Flow.

#### CODE CONFIDENTLY.

Making major changes to large codebases can be scary. Flow helps you refactor safely, so you can focus on the changes you want to make, and stop worrying about what you might break.

#### CODE BIGGER.

Working in a codebase with lots of developers can make it difficult to keep your master branch working. Flow can help prevent bad rebases. Flow can help protect your carefully designed library from your coworkers. And Flow can help you understand the code you wrote six months ago.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/dependence-nodejs-flow:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/facebook/flow.git" \
    hub.opshub.sh/containerops/dependence-nodejs-flow:latest
```

### Parameters 

Required:

- `git-url` where your code is located

### Versions 1.0.0