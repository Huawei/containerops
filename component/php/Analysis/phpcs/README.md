## PHP Code Analysis Component PHPCS

### What's the Component?

This image is php runtime image, used for analysis your php coding style. 

PHP\_CodeSniffer is a set of two PHP scripts; the main phpcs script that tokenizes PHP, JavaScript and CSS files to detect violations of a defined coding standard, and a second phpcbf script to automatically correct coding standard violations. PHP\_CodeSniffer is an essential development tool that ensures your code remains clean and consistent.

### Learn how to build it?

Use the docker build command build the image.

```shell
docker build -t hub.opshub.sh/containerops/analysis-php-phpcs:latest .
```

### Component Usage

```shell
docker run \
    --env CO_DATA=" \
    git-url=https://github.com/squizlabs/PHP_CodeSniffer.git \
    report=full standard=phpcs.xml.dist" \
    hub.opshub.sh/containerops/analysis-php-phpcs:latest
```

### Parameters 

Required:

- `git-url` where your code is located

Optional:

- `file` One or more files and/or directories to check
- `report` Print either the "full", "xml", "checkstyle", "csv","json", "junit", "emacs", "source", "summary", "diff","svnblame", "gitblame", "hgblame" or "notifysend" report (the "full" report is printed by default)
- `basepath` A path to strip from the front of file paths inside reports
- `bootstrap` A comma separated list of files to run before processing begins
- `severity` The minimum severity required to display an error or warning
- `error-severity`
- `warning-severity`
- `standard` The name or path of the coding standard to use
- `sniffs` A comma separated list of sniff codes to include or exclude from checking.(all sniffs must be part of the specified standard)
- `encoding` The encoding of the files being checked (default is utf-8)
- `parallel`
- `generator` Uses either the "HTML", "Markdown" or "Text" generator.(forces documentation generation instead of checking)
- `extensions` A comma separated list of file extensions to check.
 - (extension filtering only valid when checking a directory)
 - The type of the file can be specified using: ext/type.e.g., module/php,es/js
- `ignore`
- `file-lis`

### Versions 1.0.0