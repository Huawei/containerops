#!/bin/bash

IMAGE_NAME=`cat imagename`
docker run --rm --env CO_DATA="git-url=https://github.com/wangkirin/demo-bmi.git out-put-type=xml version=gradle3 report-path=./calculator" $IMAGE_NAME
