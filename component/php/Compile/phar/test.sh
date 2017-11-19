#!/bin/bash

function action1()
{
    export CO_DATA="git-url=https://github.com/SegmentFault/phar-sample.git entry-file=build.php"
    ./bin/containerops-php Phar
}

function action2()
{
    export CO_DATA="git-url=https://github.com/SegmentFault/phar-sample.git entry-file=build1.php"
    ./bin/containerops-php Phar
}

echo "Testing action$1"

"action$1"

rm -rf ./workspace