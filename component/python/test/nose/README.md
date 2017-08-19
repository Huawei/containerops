## BUILD

```bash
docker build -t containerops/nose .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/nose-devs/nose.git entry-path=unit_tests' containerops/nose
```

## TEST with deference python version

`version` is one of `python`, `python2`, `python3`, `py3k`.
default is `py3k`

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/nose-devs/nose.git entry-path=unit_tests version=python' containerops/nose
```
