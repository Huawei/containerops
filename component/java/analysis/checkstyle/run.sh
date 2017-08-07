#!/bin/bash

IMAGE_NAME=`cat imagename`
docker run --rm --env CO_DATA="git-url=https://githubx.com/vanniuner/gradle-demo.git out-put-type=json" $IMAGE_NAME
