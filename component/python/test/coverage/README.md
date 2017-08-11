## BUILD

```bash
docker build -t containerops/coverage .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=test/test_regex.py' containerops/coverage
```
