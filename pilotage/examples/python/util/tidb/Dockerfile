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

FROM golang:1.7

MAINTAINER xiechuan  <xiechuanj@gmail.com>

ENV HOME=/root
ENV GOPATH=$HOME/gopath \
    PATH=$HOME/gopath/bin:$PATH \
    TRAVIS_BUILD_DIR=$HOME/gopath/src/github.com/pingcap/tidb 
#    CURDIR=$HOME/gopath/src/github.com/pingcap/tidb    

#RUN git clone https://github.com/pingcap/tidb.git /root/gopath/src/github.com/pingcap/tidb

RUN apt-get update && apt-get install -y vim
RUN mkdir -p /root/gopath/src/github.com/pingcap


WORKDIR /root/gopath/src/github.com/pingcap

COPY * ./

CMD  python pingcaptidbtest.py
