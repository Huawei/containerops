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

FROM containerops/golang:1.8.1
MAINTAINER thomas.tan <tanhaijun@huawei.com>

# install python(2.7.12) pylint(1.5.2)
RUN apt-get update && apt-get install -y python python-dev python-pip python-virtualenv pylint

# bulid pylint component
RUN mkdir -p $GOPATH/src/github.com/Huawei/containerops/
ADD src/copspylint.go $GOPATH/src/github.com/Huawei/containerops/
ADD src/pylint.conf $GOPATH/src/github.com/Huawei/containerops/

ENV PATH /user/bin:$GOPATH/src/github.com/Huawei/containerops:$PATH

# defualut pylint.conf generate by cmd "pylint --generate-rcfile > pylint.conf"
# and modify set reoprts=no
ENV LINTCONFFILE $GOPATH/src/github.com/Huawei/containerops/pylint.conf

WORKDIR $GOPATH/src/github.com/Huawei/containerops
RUN go build -o copspylint copspylint.go

CMD copspylint
