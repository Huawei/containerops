## BUILD

```bash
docker build -t containerops/coverage .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=test/test_regex.py' containerops/coverage
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=test/test_regex.py version=python' containerops/coverage
```
