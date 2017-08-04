FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER vanniuner <461382141@qq.com>

RUN apt-get update
RUN apt-get install -y unzip tar openjdk-8-jdk wget git

WORKDIR /opt/gradle
RUN wget https://downloads.gradle.org/distributions/gradle-4.1-rc-1-bin.zip \
    && unzip gradle-4.1-rc-1-bin.zip \
    && rm -rf gradle-4.1-rc-1-bin.zip

ENV PATH /opt/gradle/gradle-4.1-rc-1/bin:$PATH

WORKDIR /root/convert
COPY ./JsonParse.java JsonParse.java
COPY ./build.gradle build.gradle
RUN gradle build && \
    cp build/libs/convert.jar /root/convert.jar && \
    rm -rf /root/convert

WORKDIR /root
COPY ./checkstyle.xml /root/checkstyle.xml
COPY ./compile.sh /root/compile.sh
RUN chmod 777 /root/compile.sh

CMD /root/compile.sh
