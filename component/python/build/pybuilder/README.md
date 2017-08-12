## BUILD

```bash
docker build -t containerops/pybuilder .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/blanzp/amazon_examples.git entry-path=. task=run_unit_tests' containerops/pybuilder
```
