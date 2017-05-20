FROM  phusion/baseimage:0.9.21
MAINTAINER Meaglith Ma genedna@gmail.com

ENV BAZEL_VERSION 0.4.5
ENV BAZEL_WORKSPACE /var/opt/bazel

RUN echo "deb [arch=amd64] http://storage.googleapis.com/bazel-apt stable jdk1.8" | tee /etc/apt/sources.list.d/bazel.list && curl https://bazel.build/bazel-release.pub.gpg | apt-key add -
RUN apt-get update && apt-get install -y bazel 

RUN mkdir -p /var/opt/bazel
WORKDIR /var/opt/bazel
