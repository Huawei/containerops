#!/bin/bash

export CO_DATA="git-url=http://192.168.123.201/yangkghjh/easy-php.git action=install"

echo "Testing $1"

go run component/$1.go

rm -rf ./workspace