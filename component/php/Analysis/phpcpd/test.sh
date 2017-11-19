#!/bin/bash

function phpcpd1()
{
    export CO_DATA="git-url=https://github.com/TIGERB/easy-php.git"
    ./bin/containerops-php Phpcpd
}

function phpcpd2()
{
    export CO_DATA="git_url=https://github.com/yangkghjh/containerops-php.git"
    ./bin/containerops-php Phpcpd
}

echo "Testing $1"

"phpcpd$1"

rm -rf ./workspace