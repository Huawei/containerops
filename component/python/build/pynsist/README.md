## BUILD

```bash
docker build -t containerops/pynsist .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/takluyver/pynsist.git entry-file=examples/console/installer.cfg upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/pynsist' containerops/pynsist
```
