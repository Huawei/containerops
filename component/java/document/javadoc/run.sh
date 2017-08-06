#!/bin/bash

IMAGE_NAME=`cat imagename`
docker run --rm --env CO_DATA="git-url=https://github.com/vanniuner/gradle-demo.git target=x" $IMAGE_NAME
