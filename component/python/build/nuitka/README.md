## BUILD

```bash
docker build -t containerops/nuitka .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=hub.opshub.sh/lidian/test/nuitka/v0.1' containerops/nuitka
```
