#!/bin/bash

IMAGE_NAME=`cat imagename`
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git target=https://hub.opshub.sh/binary/v1/lidian/test/binary/1.1.0/javadoc.tar" $IMAGE_NAME
