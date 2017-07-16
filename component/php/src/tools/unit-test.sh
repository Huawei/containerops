#!/bin/bash

export testspace="./testsapce"

go test -v ./test/$1

rm -rf $testspace