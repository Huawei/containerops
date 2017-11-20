#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/TIGERB/easy-php.git exclude=public"
    ./bin/containerops-php Phploc
}

function action2()
{
    export CO_DATA="git-url=https://github.com/yangkghjh/containerops-php.git path=new"
    ./bin/containerops-php Phploc
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace