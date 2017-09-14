## Python pycco Docker Component

### What's the Component?
[Pycco](https://github.com/pycco-docs/pycco) is a Python port of Docco: the original quick-and-dirty, hundred-line- long, literate-programming-style documentation generator.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/pycco .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/pycco-docs/pycco.git' containerops/pycco
```

### Parameters
- `git-url` is the source git repo url
- `out-put-type` available value: yaml,json

### Versions 1.0.0
