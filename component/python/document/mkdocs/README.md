## Python mkdocs Docker Component

### What's the Component?
[Mkdocs](https://github.com/mkdocs/mkdocs/) is Project documentation with Markdown.

### Learn how to build it?
Use the `docker build` command build the image
```bash
docker build -t containerops/mkdocs .
```

### Component Usage
```bash
docker run --rm -e CO_DATA='git-url=https://github.com/mkdocs/mkdocs.git entry-path=.' containerops/mkdocs
```

### Parameters
- `git-url` is the source git repo url
- `entry-path` is the entry path for mkdocs
- `out-put-type` available value: yaml,json

### Versions 1.0.0
