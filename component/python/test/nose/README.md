## BUILD

```bash
docker build -t containerops/nose .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/nose-devs/nose.git entry-path=unit_tests' containerops/nose
```
