# Copyright 2014 - 2017 Huawei Technologies Co., Ltd. All rights reserved.
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

FROM  phusion/baseimage:0.9.21
MAINTAINER Quanyi Ma <maquanyi@huawei.com>

ARG go_version
ENV GO_VERSION ${go_version}

COPY go${GO_VERSION}.sha256 /tmp/
RUN curl -sSL https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz -o /tmp/go.tar.gz && \
    echo "$(cat /tmp/go${GO_VERSION}.sha256)  /tmp/go.tar.gz" | sha256sum -c - && \
    tar -C /var/opt -xzf /tmp/go.tar.gz && \
    rm /tmp/go.tar.gz /tmp/go${GO_VERSION}.sha256 && \
    mkdir -p /var/opt/gopath && \
    chmod -R 777 /var/opt/gopath

ENV GOROOT /var/opt/go
ENV GOPATH /var/opt/gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH:/bin

WORKDIR $GOPATH
