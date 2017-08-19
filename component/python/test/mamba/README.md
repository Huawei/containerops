## BUILD

```bash
docker build -t containerops/mamba .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/juanAFernandez/testing-with-python.git entry-file=examples/mamba_example.py' containerops/mamba
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/juanAFernandez/testing-with-python.git entry-file=examples/mamba_example.py version=python' containerops/mamba
```
