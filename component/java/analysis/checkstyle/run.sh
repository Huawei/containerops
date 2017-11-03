#!/bin/bash

IMAGE_NAME=`cat imagename`
docker run --rm --env CO_DATA="git-url=https://github.com/wangkirin/demo-bmi.git out-put-type=json report-path=./calculator version=gradle4" $IMAGE_NAME
