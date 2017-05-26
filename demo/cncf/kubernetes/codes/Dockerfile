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

FROM containerops/golang:1.8.3
MAINTAINER Quanyi Ma <genedna@gmail.com>

RUN mkdir -p $GOPATH/github.com/Huawei/containerops
ADD ./k8s.go $GOPATH/github.com/Huawei/containerops
WORKDIR $GOPATH/github.com/Huawei/containerops
RUN go get gopkg.in/src-d/go-git.v4 && go build k8s.go
