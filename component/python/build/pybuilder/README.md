## BUILD

```bash
docker build -t containerops/pybuilder .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/blanzp/amazon_examples.git entry-path=. task=run_unit_tests' containerops/pybuilder
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/blanzp/amazon_examples.git entry-path=. task=run_unit_tests version=python' containerops/pybuilder
```
