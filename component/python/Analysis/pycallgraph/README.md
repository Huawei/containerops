## BUILD

```bash
docker build -t containerops/pycallgraph .
```

## TEST

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/python-aio-periodic.git entry-file=aio_periodic/utils.py upload=hub.opshub.sh/lidian/test/pycallgraph/v0.1' containerops/pycallgraph
```

## Warning

This pycallgraph not stable on generate callgraph.
