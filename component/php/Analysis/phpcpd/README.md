## PHP Code Analysis Component PHPcpd

### What's the Component?

This image is php runtime image, used for analysis your php coding style. 

`phpcpd` is a Copy/Paste Detector (CPD) for PHP code.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/analysis-php-phpcpd:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/TIGERB/easy-php.git" \
    hub.opshub.sh/containerops/analysis-php-phpcpd:latest
```

### Parameters 

Required:

- `git-url` where your code is located

Optional:

- `path` Files and directories to analyze
- `names` A comma-separated list of file names to check [default: ["*.php"]]
- `names-exclude` A comma-separated list of file names to exclude
- `regexps-exclude` A comma-separated list of paths regexps to exclude (example: "#var/.*_tmp#")
- `exclude` Exclude a directory from 
- `min-lines` Minimum number of identical lines [default: 5]
- `min-tokens` inimum number of identical tokens [default: 70]

### Versions 1.0.0