## BUILD

```bash
docker build -t containerops/sphinx .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-path=docs' containerops/sphinx
```
