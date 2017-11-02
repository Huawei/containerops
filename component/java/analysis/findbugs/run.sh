#!/bin/bash

IMAGE_NAME=`cat imagename`
docker run --rm --env CO_DATA="version=gradle3 git-url=https://github.com/wangkirin/demo-bmi.git out-put-type=json report-path=./calculator" $IMAGE_NAME
