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

FROM java:8

MAINTAINER xiechuan  <xiechuanj@gmail.com>

ENV SONAR_VERSION=6.1 \
    SONARQUBE_HOME=/opt/sonarqube \
    SONARSCANNER_VERSION=2.8

EXPOSE 9000

RUN set -x \

    && gpg --keyserver ha.pool.sks-keyservers.net --recv-keys F1182E81C792928921DBCAB4CFCA4A29D26468DE \
    && cd /opt \
    && curl -o sonarqube.zip -fSL https://sonarsource.bintray.com/Distribution/sonarqube/sonarqube-$SONAR_VERSION.zip \
    && curl -o sonarqube.zip.asc -fSL https://sonarsource.bintray.com/Distribution/sonarqube/sonarqube-$SONAR_VERSION.zip.asc \
    && gpg --batch --verify sonarqube.zip.asc sonarqube.zip \
    && unzip sonarqube.zip \
    && mv sonarqube-$SONAR_VERSION sonarqube \
    && rm sonarqube.zip* \
    && rm -rf $SONARQUBE_HOME/bin/* \
    && curl -o sonarscanner.zip -fSL https://sonarsource.bintray.com/Distribution/sonar-scanner-cli/sonar-scanner-$SONARSCANNER_VERSION.zip \    
    && curl -o sonar-python-plugin-1.6.jar -fSL  https://sonarsource.bintray.com/Distribution/sonar-python-plugin/sonar-python-plugin-1.6.jar \    
    && curl -o sonarscanner.zip.asc -fSL https://sonarsource.bintray.com/Distribution/sonar-scanner-cli/sonar-scanner-$SONARSCANNER_VERSION.zip.asc \    
    && curl -o sonar-python-plugin-1.6.jar.asc -fSL  https://sonarsource.bintray.com/Distribution/sonar-python-plugin/sonar-python-plugin-1.6.jar.asc \    
    && gpg --batch --verify sonarscanner.zip.asc sonarscanner.zip \
    && gpg --batch --verify sonar-python-plugin-1.6.jar.asc sonar-python-plugin-1.6.jar \
    && unzip sonarscanner.zip \
    && mv sonar-scanner-$SONARSCANNER_VERSION sonar-scanner \
    && mv sonar-python-plugin-1.6.jar /opt/sonarqube/extensions/plugins/ \
    && rm sonarscanner.zip*

WORKDIR $SONARQUBE_HOME
COPY run.sh $SONARQUBE_HOME/bin/
COPY checkCode.py $SONARQUBE_HOME/
RUN chmod +x ./bin/run.sh
RUN chmod +x ./checkCode.py
ENTRYPOINT python checkCode.py
