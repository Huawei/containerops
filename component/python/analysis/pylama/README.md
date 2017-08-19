## BUILD

```bash
docker build -t containerops/pylama .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git' containerops/pylama
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git version=python' containerops/pylama
```
