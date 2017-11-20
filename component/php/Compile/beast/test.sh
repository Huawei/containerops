#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/TIGERB/easy-php.git"
    ./bin/containerops-php Beast
}

function action2()
{
    export CO_DATA="git_url=https://github.com/TIGERB/easy-php.git"
    ./bin/containerops-php Beast
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace