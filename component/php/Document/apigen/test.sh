#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/TIGERB/easy-php.git path=app destination=docs"
    ./bin/containerops-php Apigen
}

function action2()
{
    export CO_DATA="git-url=https://github.com/TIGERB/easy-php.git path=app2 destination=docs"
    ./bin/containerops-php Apigen
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace