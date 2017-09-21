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

FROM hub.opshub.sh/containerops/php:latest
MAINTAINER Yang Bin <yangkghjh@gmail.com>

USER root

RUN apt-get update && apt-get install -y libxslt-dev
ADD xsl.so /usr/local/lib/php/extensions/no-debug-non-zts-20160303/xsl.so
ADD php.ini /usr/local/etc/php/php.ini
RUN composer global require "theseer/phpdox";
ENV PATH ~/.composer/vendor/bin:$PATH
RUN apt-get install -y php-xsl

RUN mkdir -p /root/src
ADD / /root/src

WORKDIR /root/src

CMD ./bin/containerops-php Phpdox