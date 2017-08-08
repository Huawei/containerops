## BUILD

```bash
docker build -t containerops/pdoc .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-mod=grapy' containerops/pdoc
```
