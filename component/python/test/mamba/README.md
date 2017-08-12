## BUILD

```bash
docker build -t containerops/mamba .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/juanAFernandez/testing-with-python.git entry-file=examples/mamba_example.py' containerops/mamba
```
