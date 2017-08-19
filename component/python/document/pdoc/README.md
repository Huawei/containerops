## BUILD

```bash
docker build -t containerops/pdoc .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-mod=grapy' containerops/pdoc
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-mod=grapy version=python' containerops/pdoc
```
