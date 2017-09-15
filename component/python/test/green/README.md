## Python green Docker Component

### What's the Component?
[Green](https://github.com/CleanCut/green) is a clean, colorful, fast python test runner.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/green .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=.' containerops/green
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=. version=python' containerops/green
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-path` is the entry file or path for coverage
- `out-put-type` available value: yaml,json

### Versions 1.0.0
