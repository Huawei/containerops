#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/TIGERB/easy-php.git path=docs"
    ./bin/containerops-php Phpmd
}

function action2()
{
    export CO_DATA="git-url=https://github.com/yangkghjh/containerops-php.git"
    ./bin/containerops-php Phpmd
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace