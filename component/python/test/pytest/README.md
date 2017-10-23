## Python pytest Docker Component

### What's the Component?
The [pytest](https://github.com/pytest-dev/pytest/) framework makes it easy to write small tests, yet scales to support complex functional testing for applications and libraries.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/pytest .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=.' containerops/pytest
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=. version=python' containerops/pytest
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-path` is the entry file or path for pytest
- `out-put-type` available value: yaml,xml

### Versions 1.0.0
