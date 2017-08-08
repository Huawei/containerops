## BUILD

```bash
docker build -t containerops/pylint .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git' containerops/pylint
```
