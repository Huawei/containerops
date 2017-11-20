#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/sebastianbergmann/phploc.git --configuration=phpunit.xml composer=true"
    ./bin/containerops-php Phpunit
}

function action2()
{
    export CO_DATA="git-url=https://github.com/yangkghjh/containerops-php.git --configuration=phpunit1.xml composer=true"
    ./bin/containerops-php Phpunit
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace