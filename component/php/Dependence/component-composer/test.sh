#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/sebastianbergmann/phploc.git"
    ./bin/containerops-php Composer
}

function action2()
{
    export CO_DATA="git-url=https://github.com/yangkghjh/containerops-php.git"
    ./bin/containerops-php Composer
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace