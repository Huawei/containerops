#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/wp-cli/wp-cli.git file=./bin/wp composer=true"
    ./bin/containerops-php Cli
}

function action2()
{
    export CO_DATA="git-url=https://github.com/yangkghjh/containerops-php.git file=./bin/wp composer=false"
    ./bin/containerops-php Cli
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace