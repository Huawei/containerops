## BUILD

```bash
docker build -t containerops/dh-virtualenv .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/spotify/dh-virtualenv.git upload=hub.opshub.sh/lidian/test/nuitka/v0.1' containerops/dh-virtualenv
```
