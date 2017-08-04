FROM huawei

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
