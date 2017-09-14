## Python pybuilder Docker Component

### What's the Component?
[PyBuilder](https://github.com/pybuilder/pybuilder) is a software build tool written in 100% pure Python, mainly targeting Python applications.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/pybuilder .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/blanzp/amazon_examples.git entry-path=. task=run_unit_tests' containerops/pybuilder
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/blanzp/amazon_examples.git entry-path=. task=run_unit_tests version=python' containerops/pybuilder
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-path` is the entry path with `build.py` for pybuilder
- `task` is the task name of pybuilder
- `out-put-type` available value: yaml,json

### Versions 1.0.0
