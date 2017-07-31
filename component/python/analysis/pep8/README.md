## BUILD

```bash
docker build -t containerops/pep8 .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git' containerops/pep8
```
