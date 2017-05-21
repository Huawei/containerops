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

FROM ubuntu:14.04

MAINTAINER xiechuan  <xiechuanj@gmail.com>

RUN apt-get update && apt-get install -y --no-install-recommends \
    curl \
    file \
    make \
    build-essential \
    checkinstall \
    libgflags-dev \
    libsnappy-dev \
    liblz4-dev \
    zlib1g-dev \
    libbz2-dev \
    ca-certificates \
    git \
    python \
    && rm -rf /var/lib/apt/lists/* && apt-get clean

RUN cd / && \
    curl -L https://github.com/facebook/rocksdb/archive/4.12.fb.tar.gz -o rocksdb.tar.gz && \
    tar xzf rocksdb.tar.gz && \
    cd rocksdb-4.12.fb && \
    PORTABLE=1 make shared_lib && \
    make install && \
    cd / && \
    rm -rf /rocksdb-4.12.fb /rocksdb.tar.gz

RUN curl -sSf https://static.rust-lang.org/rustup.sh | sh -s  -- --date=2016-10-06 --disable-sudo -y --channel=nightly


COPY *.py ./

CMD python pingcaptikvtest.py

