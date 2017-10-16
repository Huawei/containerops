# PHPCS

## Build

```shell
docker build -t hub.opshub.sh/containerops/analysis-php-phpcs:latest .
```

## Run

```shell
docker run --env CO_DATA="git-url=https://github.com/squizlabs/PHP_CodeSniffer.git report=full standard=phpcs.xml.dist" hub.opshub.sh/containerops/analysis-php-phpcs:latest
```

## Options

Required:

- git-url

Optional:

- file
- report
- basepath
- bootstrap
- severity
- error-severity
- warning-severity
- standard
- sniffs
- encoding
- parallel
- generator
- extensions
- ignore
- file-lis

```shell
<basepath>     A path to strip from the front of file paths inside reports
<bootstrap>    A comma separated list of files to run before processing begins
<file>         One or more files and/or directories to check
<encoding>     The encoding of the files being checked (default is utf-8)
<extensions>   A comma separated list of file extensions to check
            (extension filtering only valid when checking a directory)
            The type of the file can be specified using: ext/type
            e.g., module/php,es/js
<generator>    Uses either the "HTML", "Markdown" or "Text" generator
            (forces documentation generation instead of checking)
<report>       Print either the "full", "xml", "checkstyle", "csv"
            "json", "junit", "emacs", "source", "summary", "diff"
            "svnblame", "gitblame", "hgblame" or "notifysend" report
            (the "full" report is printed by default)
<severity>     The minimum severity required to display an error or warning
<sniffs>       A comma separated list of sniff codes to include or exclude from checking
            (all sniffs must be part of the specified standard)
<standard>     The name or path of the coding standard to use
```