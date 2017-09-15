## Python sphinx Docker Component

### What's the Component?
[Sphinx](https://github.com/sphinx-doc/sphinx/) is a tool that makes it easy to create intelligent and beautiful documentation, written by Georg Brandl and licensed under the BSD license.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/sphinx .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-path=docs' containerops/sphinx
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/grapy.git entry-path=docs version=python' containerops/sphinx
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-path` is the entry document path for sphinx
- `out-put-type` available value: yaml,json

### Versions 1.0.0
