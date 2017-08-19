## BUILD

```bash
docker build -t containerops/tox .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/CleanCut/green.git entry-path=.' containerops/tox
```
