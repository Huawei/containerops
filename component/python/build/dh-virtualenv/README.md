## BUILD

```bash
docker build -t containerops/dh-virtualenv .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/spotify/dh-virtualenv.git upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/dh-virtualenv' containerops/dh-virtualenv
```
