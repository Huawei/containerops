FROM java:8

MAINTAINER xiechuan  <xiechuanj@gmail.com>

#ENV http_proxy=http://web-proxy.rose.hp.com:8080 \
#    https_proxy=http://web-proxy.rose.hp.com:8080 \
#    HTTP_PROXY=http://web-proxy.rose.hp.com:8080 \
#    HTTPS_PROXY=http://web-proxy.rose.hp.com:8080 \
ENV SONAR_VERSION=6.1 \
    SONARQUBE_HOME=/opt/sonarqube \
    SONARSCANNER_VERSION=2.8
#    CO_DATA='{"COMPONENT_START": "http://localhost:8000",\n"COMPONENT_STOP": "http://localhost:8000",\n"TASK_START": "http://localhost:8000",\n"TASK_RESULT": "http://172.17.0.1:8001",\n"TASK_STATUS": "http://localhost:8000","REGISTER_URL": "http://172.17.0.1:8001",\n"SERVICE_ADDR": "127.0.0.1:30001:9999",\n"RUN_ID": "1",\n"POD_NAME":"go-doc1-2-1-178-1-pod",\n"gitUrl": "https://github.com/xiechuanj/python-sonar-runner.git",\n"contents": "sonar.projectKey=python-sonar-runner\nsonar.projectName=python-sonar-runner\nsonar.projectVersion=1.0\nsonar.sources=src\nsonar.language=py\nsonar.sourceEncoding=UTF-8",\n"filename": "sonar-project.properties",\n"EVENT_LIST": "COMPONENT_START,1;COMPONENT_STOP,2;TASK_START,8;TASK_RESULT,9;TASK_STATUS,10;REGISTER_URL,11"}'

# Http port
EXPOSE 9000

RUN set -x \

    # pub   2048R/D26468DE 2015-05-25
    #       Key fingerprint = F118 2E81 C792 9289 21DB  CAB4 CFCA 4A29 D264 68DE
    # uid                  sonarsource_deployer (Sonarsource Deployer) <infra@sonarsource.com>
    # sub   2048R/06855C1D 2015-05-25
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
