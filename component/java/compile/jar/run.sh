#!/bin/bash

IMAGE_NAME=`cat imagename`
docker run --rm --env CO_DATA="git-url=https://github.com/wangkirin/demo-bmi.git target=https://hub.opshub.sh/binary/v1/containerops/component/binary/2.2.4/demo.jar version=gradle3 build-path=./calculator" $IMAGE_NAME