## Python nose Docker Component

### What's the Component?
[nose](https://github.com/nose-devs/nose) is nicer testing for python

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/nose .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/nose-devs/nose.git entry-path=unit_tests' containerops/nose
docker run --rm -e CO_DATA='git-url=https://github.com/nose-devs/nose.git entry-path=unit_tests version=python' containerops/nose
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-path` is the entry file or path for nose
- `out-put-type` available value: yaml,xml

### Versions 1.0.0
