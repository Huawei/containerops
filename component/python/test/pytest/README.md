## BUILD

```bash
docker build -t containerops/pytest .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=.' containerops/pytest
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=. version=python' containerops/pytest
```
