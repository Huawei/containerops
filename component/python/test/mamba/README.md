## Python mamba Docker Component

### What's the Component?
[mamba](https://github.com/nestorsalceda/mamba) is the definitive test runner for Python. Born under the banner of behavior-driven development.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/mamba .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/juanAFernandez/testing-with-python.git entry-file=examples/mamba_example.py' containerops/mamba
docker run --rm -e CO_DATA='git-url=https://github.com/juanAFernandez/testing-with-python.git entry-file=examples/mamba_example.py version=python' containerops/mamba
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-file` is the entry file for mamba
- `out-put-type` available value: yaml,json

### Versions 1.0.0
