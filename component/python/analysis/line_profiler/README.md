## BUILD

```bash
docker build -t containerops/line_profiler .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/istrategylabs/python-profiling entry-file=debug.py' containerops/line_profiler
```
