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
MAINTAINER Zhen Ju <juzhenatpku@gmail.com>

USER root

RUN apt-get update && apt-get install -y gcc make g++

ENV PATH $PATH:$GOPATH/src/github.com/Huawei/containerops

RUN mkdir -p $GOPATH/src/github.com/Huawei/

WORKDIR $GOPATH/src/github.com/Huawei/
RUN git clone https://github.com/Huawei/containerops

# The essential dependencies for singular
RUN go get "github.com/cloudflare/cfssl/cli"
RUN go get "github.com/cloudflare/cfssl/cli/genkey"
RUN go get "github.com/cloudflare/cfssl/cli/sign"
RUN go get "github.com/cloudflare/cfssl/csr"
RUN go get "github.com/cloudflare/cfssl/initca"
RUN go get "github.com/cloudflare/cfssl/signer"
RUN go get "github.com/digitalocean/godo"
RUN go get "github.com/fernet/fernet-go"
RUN go get "github.com/logrusorgru/aurora"
RUN go get "github.com/mitchellh/go-homedir"
RUN go get "github.com/pkg/sftp"
RUN go get "github.com/spf13/cobra"
RUN go get "github.com/spf13/viper"
RUN go get "golang.org/x/crypto/ssh"
RUN go get "golang.org/x/net/context"
RUN go get "golang.org/x/oauth2"
RUN go get "gopkg.in/yaml.v2"

WORKDIR $GOPATH/src/github.com/Huawei/containerops/singular/
RUN go build

#ADD codes/*.go  $GOPATH/src/github.com/Huawei/containerops/demo/cncf/
#ADD singular.template.yaml $GOPATH/src/github.com/Huawei/containerops/demo/cncf/

WORKDIR $GOPATH/src/github.com/Huawei/containerops/demo/cncf/singular/codes/
RUN go build -o ../singular-demo

WORKDIR $GOPATH/src/github.com/Huawei/containerops/demo/cncf/singular/
RUN cp $GOPATH/src/github.com/Huawei/containerops/singular/singular ./

CMD ./singular-demo
