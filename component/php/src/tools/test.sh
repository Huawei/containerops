#!/bin/bash

function phar()
{
    export CO_DATA="git-url=https://github.com/SegmentFault/phar-sample.git entry-file=build.php"
    go run component/phar.go
}

function phpunit()
{
    export CO_DATA="git-url=https://github.com/sebastianbergmann/phploc.git --configuration=phpunit.xml composer=true"
    go run component/phpunit.go
}

echo "Testing $1"

"$1"

rm -rf ./workspace