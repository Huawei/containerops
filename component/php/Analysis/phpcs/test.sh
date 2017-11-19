#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/squizlabs/PHP_CodeSniffer.git report=full standard=phpcs.xml.dist"
    ./bin/containerops-php Phpcs
}

function action2()
{
    export CO_DATA="git-url=https://github.com/yangkghjh/containerops-php.git report=full standard=phpcs.xml.dist"
    ./bin/containerops-php Phpcs
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace