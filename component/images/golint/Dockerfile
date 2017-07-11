# Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# base image: ubuntu uname:3.13.0-24-generic #47-Ubuntu SMP Fri May 2 23:30:00 UTC 2014 x86_64 x86_64 x86_64 GNU/Linux
FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER thomas.tan <tanhaijun@huawei.com>

RUN apt-get update && apt-get install -y git zip unzip
RUN curl -sSL https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz -o /tmp/go.tar.gz && \
  echo "1862f4c3d3907e59b04a757cfda0ea7aa9ef39274af99a784f5be843c80c6772  /tmp/go.tar.gz" | sha256sum -c - && \
  tar -C /var/opt -xzf /tmp/go.tar.gz && \
  rm /tmp/go.tar.gz && \
  mkdir -p /var/opt/gopath && \
  chmod -R 777 /var/opt/gopath



ENV GOROOT /var/opt/go
ENV GOPATH /var/opt/gopath

# install gotools
RUN curl -sSL https://github.com/golang/tools/archive/master.zip -o /tmp/gotool.zip && \
  unzip -o -d  /tmp /tmp/gotool.zip && \
  mkdir -p $GOPATH/src/golang.org/x/tools && \
  cp -r /tmp/tools-master/*  $GOPATH/src/golang.org/x/tools && \
  rm  /tmp/gotool.zip && \
  rm /tmp/tools-master/ -r


# install golint
RUN curl -sSL https://github.com/golang/lint/archive/master.zip -o /tmp/golint.zip && \
  unzip -o -d  /tmp /tmp/golint.zip && \
  mkdir -p $GOPATH/src/github.com/golang/lint && \
  cp -r /tmp/lint-master/*  $GOPATH/src/github.com/golang/lint && \
  rm /tmp/lint-master/ -r  && \
  rm /tmp/golint.zip

ENV PATH $GOPATH/src/github.com/golang/lint/golint:$GOROOT/bin:$GOPATH/bin:$GOPATH/src/github.com/Huawei/containerops/:$PATH

WORKDIR $GOPATH/src/github.com/golang/lint/golint
RUN go build ./


# start golint component
RUN mkdir -p $GOPATH/src/github.com/Huawei/containerops/
ADD src/copsgolint.go $GOPATH/src/github.com/Huawei/containerops/

WORKDIR $GOPATH/src/github.com/Huawei/containerops
RUN go build -o copsgolint copsgolint.go

# WORKDIR $GOPATH/src/github.com/golang/lint/golint

CMD copsgolint
