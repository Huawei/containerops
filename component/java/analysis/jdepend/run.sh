#!/bin/bash

IMAGE_NAME=`cat imagename`
docker run --rm --env CO_DATA="git-url=https://github.com/wangkirin/demo-bmi.gitout-put-type=json version=gradle3 report-path=./calculator" $IMAGE_NAME
