FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER fred.liu <461382141@qq.com>

RUN apt-get update
RUN apt-get install -y unzip tar openjdk-8-jdk wget git

WORKDIR /opt/gradle
RUN wget https://services.gradle.org/distributions/gradle-3.5-bin.zip \
    && unzip gradle-3.5-bin.zip \
    && rm -rf gradle-3.5-bin.zip

ENV PATH /opt/gradle/gradle-3.5/bin:$PATH

WORKDIR /root
COPY ./compile.sh /root/compile.sh
RUN chmod 777 /root/compile.sh

CMD /root/compile.sh