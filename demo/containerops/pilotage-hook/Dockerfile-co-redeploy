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

FROM hub.opshub.sh/containerops/golang:1.8.3
MAINTAINER Zhen Ju<juzhen@huawei.com>

ADD ./codes/co_redeploy.go /
ADD ./codes/common.go /

RUN go get -v golang.org/x/crypto/ssh
WORKDIR /
RUN go build -o /usr/local/bin/coredeploy common.go co_redeploy.go
RUN mkdir -p $GOPATH/src/github.com/Huawei/
CMD coredeploy
