## Node.js Code Document Component DOCCO

### What's the Component?

This image is node.js runtime image, used to generate document for your node.js project.

Docco is a quick-and-dirty documentation generator, written in Literate CoffeeScript. It produces an HTML document that displays your comments intermingled with your code. All prose is passed through Markdown, and code is passed through Highlight.js syntax highlighting. This page is the result of running Docco against its own source file.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/document-nodejs-docco:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git_url=https://github.com/jonpliske/coffeescript_examples.git \
    file=*.coffee" \
    hub.opshub.sh/containerops/document-nodejs-docco:latest
```

### Parameters 

Required:

- `git-url` where your code is located
- `file` targe files

### Versions 1.0.0
