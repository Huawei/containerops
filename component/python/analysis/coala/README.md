## Python Coala code analysis Docker Component

### What's the Component?
[coala](https://github.com/coala/coala) provides a unified command-line interface for linting and fixing all your code,
regardless of the programming languages you use.

### Learn how to build it?
Use the `docker build` command build the image

```bash
docker build -t containerops/coala .
```

### Component Usage

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git' containerops/coala
```

### Parameters

- `git-url` is the source git repo url

### Versions 1.0.0
