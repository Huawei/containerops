## BUILD

```bash
docker build -t containerops/pytest .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=.' containerops/pytest
```
