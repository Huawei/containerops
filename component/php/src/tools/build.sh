#!/bin/bash

echo "Start building...\n"

case $1 in
"phar")
    docker build -t containerops/phar:latest -f Compile/phar/Dockerfile .
    ;;
"phpunit")
    docker build -t containerops/phpunit:latest -f Unittest/phpunit/Dockerfile .
    ;;
*)
    echo "No such component: $1.\n"
    exit
    ;;
esac

echo "Build success."
