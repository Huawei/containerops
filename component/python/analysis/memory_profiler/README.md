## BUILD

```bash
docker build -t containerops/memory_profiler .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/fabianp/memory_profiler.git entry-file=test/test_func.py' containerops/memory_profiler
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/fabianp/memory_profiler.git entry-file=test/test_func.py version=python' containerops/memory_profiler
```
