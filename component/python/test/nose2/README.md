## BUILD

```bash
docker build -t containerops/nose2 .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=.' containerops/nose2
```
