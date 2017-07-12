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

FROM containerops/composer:latest
MAINTAINER Yang Bin <yangkghjh@gmail.com>

RUN useradd --create-home --no-log-init --shell /bin/bash composer
RUN adduser composer sudoRUN echo 'composer:composer' | chpasswd
USER composer

RUN mkdir -p $GOPATH/src/github.com/Huawei/containerops
ADD /src/ $GOPATH/src/github.com/Huawei/containerops/

WORKDIR $GOPATH/src/github.com/Huawei/containerops
RUN go build component/component-composer.go

CMD component-composer