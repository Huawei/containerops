## Python nuitka Docker Component

### What's the Component?
[Nuitka](https://github.com/kayhayen/Nuitka) is the Python compiler.
It is a seamless replacement or extension to the Python interpreter and compiles every construct that CPython 2.6, 2.7, 3.2, 3.3, 3.4, 3.5, and 3.6 have.
It then executes uncompiled code, and compiled code together in an extremely compatible manner.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/nuitka .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/nuitka' containerops/nuitka
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/nuitka version=python' containerops/nuitka
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-file` is the entry file for nuitka
- `upload` is the upload url with PUT method for build result

### Versions 1.0.0
