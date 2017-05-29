# Copyright 20146 - 2017 Huawei Technologies Co., Ltd. All rights reserved.
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

FROM golang
MAINTAINER Marx Guo "go.morefork@gmail.com"
RUN "mkdir" "-p" "/go/src/github.com/kubernetes"
WORKDIR /go/src/github.com/kubernetes
RUN "git" "clone" "https://github.com/kubernetes/kubernetes.git"
ADD ./generateGodoc.go /go/src/
WORKDIR /go/src
CMD go run generateGodoc.go
