## BUILD

```bash
docker build -t containerops/pynsist .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/takluyver/pynsist.git entry-file=examples/console/installer.cfg upload=hub.opshub.sh/lidian/test/pynsist/v0.1' containerops/pynsist
```
