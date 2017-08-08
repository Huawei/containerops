## BUILD

```bash
docker build -t containerops/mkdocs .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/mkdocs/mkdocs.git entry-path=.' containerops/mkdocs
```
