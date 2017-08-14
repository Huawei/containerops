## BUILD

```bash
docker build -t containerops/pyinstaller .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=hub.opshub.sh/lidian/test/pyinstaller/v0.1' containerops/pyinstaller
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=hub.opshub.sh/lidian/test/pyinstaller/v0.1 version=python' containerops/pyinstaller
```
