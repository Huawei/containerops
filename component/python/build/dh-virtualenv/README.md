## Python dh-virtualenv Docker Component

### What's the Component?
[dh-virtualenv](https://github.com/spotify/dh-virtualenv) is a tool that aims to combine Debian packaging with self-contained virtualenv based Python deployments.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/dh-virtualenv .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/spotify/dh-virtualenv.git upload=https://hub.opshub.sh/binary/v1/containerops/component/binary/v0.1/dh-virtualenv' containerops/dh-virtualenv
```

### Parameters
- `git-url` is the source git repo url
- `upload` is the upload url with PUT method for build result

### Versions 1.0.0
