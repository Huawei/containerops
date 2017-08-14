## BUILD

```bash
docker build -t containerops/line_profiler .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/istrategylabs/python-profiling entry-file=debug.py' containerops/line_profiler
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/istrategylabs/python-profiling entry-file=debug.py version=python' containerops/line_profiler
```
