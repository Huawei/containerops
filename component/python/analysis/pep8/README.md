## Python pep8 code analysis Docker Component

### What's the Component?
[pep8](https://github.com/PyCQA/pycodestyle) is a tool to check your Python code against some of the style conventions in PEP 8.

### Learn how to build it?
Use the `docker build` command build the image

```bash
docker build -t containerops/pep8 .
```

### Component Usage

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git' containerops/pep8
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git version=python' containerops/pep8
```

### Parameters

- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `out-put-type` available value: yaml,json

### Versions 1.0.0
