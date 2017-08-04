#!/bin/bash
docker stop $(docker ps | awk '{ print $1 }')
docker rm $(docker ps -a | awk '{ print $1 }')
docker rmi $(docker images | awk '/^<none>/ { print $3 }')
docker rmi $(docker images | awk '/^`cat imagename` { print $3 }')
