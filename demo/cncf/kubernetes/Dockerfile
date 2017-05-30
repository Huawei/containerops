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

FROM containerops/bazel:latest
MAINTAINER Quanyi Ma <genedna@gmail.com>

USER root

RUN apt-get update && apt-get install -y git make python python-dev python-pip python-virtualenv
RUN curl -sSL https://storage.googleapis.com/golang/go1.8.3.linux-amd64.tar.gz -o /tmp/go.tar.gz && \
  echo "1862f4c3d3907e59b04a757cfda0ea7aa9ef39274af99a784f5be843c80c6772  /tmp/go.tar.gz" | sha256sum -c - && \
  tar -C /var/opt -xzf /tmp/go.tar.gz && \
  rm /tmp/go.tar.gz && \
  mkdir -p /var/opt/gopath && \
  chmod -R 777 /var/opt/gopath

ENV GOROOT /var/opt/go
ENV GOPATH /var/opt/gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH:/bin:$GOPATH/src/github.com/Huawei/containerops

RUN mkdir -p $GOPATH/src/github.com/Huawei/containerops
ADD codes/*.go $GOPATH/src/github.com/Huawei/containerops/

WORKDIR $GOPATH/src/github.com/Huawei/containerops

RUN go build k8s.go

WORKDIR $GOPATH/src/github.com/kubernetes/kubernetes

CMD k8s
