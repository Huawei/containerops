## Python line_profiler code analysis Docker Component

### What's the Component?
[line_profiler](https://github.com/rkern/line_profiler) is a module for doing line-by-line profiling of functions.

### Learn how to build it?
Use the `docker build` command build the image

```bash
docker build -t containerops/line_profiler .
```

### Component Usage

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/istrategylabs/python-profiling entry-file=debug.py' containerops/line_profiler
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/istrategylabs/python-profiling entry-file=debug.py version=python' containerops/line_profiler
```

### Parameters

- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-file` is the entry file for line profile
- `out-put-type` available value: yaml,json

### Versions 1.0.0
