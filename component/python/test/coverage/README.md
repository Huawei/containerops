## Python coverage Docker Component

### What's the Component?
[Coverage.py](https://github.com/nedbat/coveragepy) measures code coverage, typically during test execution. It uses the code analysis tools and tracing hooks provided in the Python standard library to determine which lines are executable, and which have been executed.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/coverage .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=test/test_regex.py' containerops/coverage
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=test/test_regex.py version=python' containerops/coverage
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-path` is the entry file or path for coverage
- `out-put-type` available value: yaml,xml

### Versions 1.0.0
