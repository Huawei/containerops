## Python pynsist Docker Component

### What's the Component?
[pynsist](https://github.com/takluyver/pynsist) is a tool to build Windows installers for your Python applications. The installers bundle Python itself, so you can distribute your application to people who don't have Python installed.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/pynsist .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/takluyver/pynsist.git entry-file=examples/console/installer.cfg upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/pynsist' containerops/pynsist
```

### Parameters
- `git-url` is the source git repo url
- `entry-file` is the config file of pynsist
- `upload` is the upload url with PUT method for build result

### Versions 1.0.0
