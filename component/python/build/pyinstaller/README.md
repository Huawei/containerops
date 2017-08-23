## Python pyinstaller Docker Component

### What's the Component?
[PyInstaller](https://github.com/pyinstaller/pyinstaller) bundles a Python application and all its dependencies into a single package. The user can run the packaged app without installing a Python interpreter or any modules.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/pyinstaller .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/pyinstaller' containerops/pyinstaller
# test with python2
docker run --rm -e CO_DATA='git-url=https://github.com/Lupino/bpnn.git entry-file=bpnn.py upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/pyinstaller version=python' containerops/pyinstaller
```

### Parameters
- `git-url` is the source git repo url
- `version` is one of `python`, `python2`, `python3`, `py3k`.  default is `py3k`
- `entry-file` is the entry file for pyinstaller
- `upload` is the upload url with PUT method for build result

### Versions 1.0.0
