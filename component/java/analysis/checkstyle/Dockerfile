FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER fred.liu <461382141@qq.com>

RUN apt-get update && apt-get install -y unzip tar openjdk-8-jdk wget git

WORKDIR /opt/gradle
RUN wget https://services.gradle.org/distributions/gradle-4.0-bin.zip \
    && rm -rf gradle-4.0-bin \
    && unzip gradle-4.0-bin.zip \
    && rm -rf gradle-4.0-bin.zip
RUN wget https://services.gradle.org/distributions/gradle-3.5-bin.zip \
    && rm -rf gradle-3.5-bin \
    && unzip gradle-3.5-bin.zip \
    && rm -rf gradle-3.5-bin.zip

ENV gradle3 /opt/gradle/gradle-3.5/bin
ENV gradle4 /opt/gradle/gradle-4.0/bin
ENV PATH /opt/gradle/gradle-4.0/bin:/opt/gradle/gradle-3.5/bin:$PATH

WORKDIR /root
COPY ./checkstyle.xml /root/checkstyle.xml
COPY ./checkstyle.conf /root/checkstyle.conf
COPY ./compile.sh /root/compile.sh
RUN chmod 777 /root/compile.sh

WORKDIR /root/convert
COPY ./Convert.java Convert.java
COPY ./build.gradle build.gradle
RUN mkdir -p ./config/checkstyle && \
    cp /root/checkstyle.xml ./config/checkstyle/ && \
    $gradle3/gradle build && \
    cp build/libs/convert.jar /root/convert.jar && \
    rm -rf /root/convert

CMD /root/compile.sh