## Python pylint code analysis Docker Component

### What's the Component?
[Pylint](https://github.com/PyCQA/pylint) is a Python source code analyzer which looks for programming errors, helps enforcing a coding standard and sniffs for some code smells (as defined in Martin Fowler's Refactoring book).

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/pylint .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git' containerops/pylint
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git version=python' containerops/pylint
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `out-put-type` available value: yaml,json

### Versions 1.0.0
