#!/bin/bash

export CO_DATA="git-url=http://192.168.123.201/yangkghjh/PHP_CodeSniffer.git report=full standard=phpcs.xml.dist"

echo "Testing $1"

go run component/$1.go

rm -rf ./workspace