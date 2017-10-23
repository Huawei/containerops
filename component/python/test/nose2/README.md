## Python nose2 Docker Component

### What's the Component?
[nose2](https://github.com/nose-devs/nose2) is the next generation of nicer testing for Python, based on the plugins branch of unittest2.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/nose2 .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=.' containerops/nose2
docker run --rm -e CO_DATA='git-url=https://github.com/minhhh/regex.git entry-path=. version=python' containerops/nose2
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-path` is the entry file or path for nose2
- `out-put-type` available value: yaml,xml

### Versions 1.0.0
