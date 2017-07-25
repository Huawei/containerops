#!/bin/bash

function phar()
{
    docker run --env CO_DATA="git-url=https://github.com/SegmentFault/phar-sample.git entry-file=build.php" containerops/phar:latest
}

echo "Run $1"

"$1"

rm -rf ./workspace