## Python pdoc Docker Component

### What's the Component?
[pdoc](https://github.com/BurntSushi/pdoc) is a simple command line tool and library to auto generate API documentation for Python libraries.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/pdoc .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-mod=grapy' containerops/pdoc
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-mod=grapy version=python' containerops/pdoc
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-mod` is the entry module name for pdoc you want to document
- `out-put-type` available value: yaml,json

### Versions 1.0.0
