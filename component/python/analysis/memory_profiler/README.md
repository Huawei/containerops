## BUILD

```bash
docker build -t containerops/memory_profiler .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/fabianp/memory_profiler.git entry-file=test/test_func.py' containerops/memory_profiler
```
