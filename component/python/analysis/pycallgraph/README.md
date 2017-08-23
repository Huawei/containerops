## Python pycallgraph code analysis Docker Component

### What's the Component?
[Python Call Graph](https://github.com/gak/pycallgraph) is a Python module that creates call graph visualizations for Python applications.

### Learn how to build it?
Use the `docker build` command build the image

```bash
docker build -t containerops/pycallgraph .
```

### Component Usage

```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/pycallgraph' containerops/pycallgraph
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/pycallgraph version=python' containerops/pycallgraph
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-file` is the entry file for pycallgraph
- `upload` is the output image upload url with PUT method

### Versions 1.0.0
