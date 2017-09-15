## Python tox Docker Component

### What's the Component?

[Tox](https://github.com/tox-dev/tox) is a generic virtualenv management and test command line tool you can use for:

- checking your package installs correctly with different Python versions and
- interpreters
- running your tests in each of the environments, configuring your test tool of choice
- acting as a frontend to Continuous Integration servers, greatly reducing boilerplate and merging CI and shell-based testing.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/tox .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/CleanCut/green.git entry-path=.' containerops/tox
```

### Parameters
- `git-url` is the source git repo url
- `entry-path` is the entry path for tox
- `out-put-type` available value: yaml,json

### Versions 1.0.0
