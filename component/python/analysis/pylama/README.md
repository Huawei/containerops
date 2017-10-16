## Python pylama code analysis Docker Component

### What's the Component?
[pylama](https://github.com/klen/pylama) is code audit tool for python.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/pylama .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git' containerops/pylama
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git version=python' containerops/pylama
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `out-put-type` available value: yaml,json

### Versions 1.0.0
