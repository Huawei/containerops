#!/bin/bash

IMAGE_NAME=`cat imagename`
docker build -t $IMAGE_NAME ./
