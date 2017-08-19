## BUILD

```bash
docker build -t containerops/pycco .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/pycco-docs/pycco.git' containerops/pycco
```
