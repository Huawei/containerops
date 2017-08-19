## BUILD

```bash
docker build -t containerops/sphinx .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-path=docs' containerops/sphinx
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-path=docs version=python' containerops/sphinx
```
