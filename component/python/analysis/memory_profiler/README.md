## Python memory_profiler code analysis Docker Component

### What's the Component?
[memory_profiler](https://github.com/fabianp/memory_profiler) is a python module for monitoring memory consumption of a process as well as line-by-line analysis of memory consumption for python programs.

### Learn how to build it?
Use the `docker build` command build the image

```bash
docker build -t containerops/memory_profiler .
```

### Component Usage

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/fabianp/memory_profiler.git entry-file=test/test_func.py' containerops/memory_profiler
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/fabianp/memory_profiler.git entry-file=test/test_func.py version=python' containerops/memory_profiler
```

### Parameters

- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-file` is the entry file for memory profile
- `out-put-type` available value: yaml,json

### Versions 1.0.0
