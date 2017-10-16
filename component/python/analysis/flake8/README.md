## Python Flake8 code analysis Docker Component

### What's the Component?
[Flake8](https://github.com/PyCQA/flake8) is a wrapper around these tools:

* PyFlakes
* pycodestyle
* Ned Batchelder's McCabe script

Flake8 runs all the tools by launching the single flake8 command.
It displays the warnings in a per-file, merged output.

### Learn how to build it?
Use the `docker build` command build the image

```bash
docker build -t containerops/flake8 .
```

### Component Usage

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git' containerops/flake8
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git version=python' containerops/flake8
```

### Parameters

- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`. default is `py3k`
- `out-put-type` available value: yaml,json

### Versions 1.0.0
