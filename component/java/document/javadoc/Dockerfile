FROM huawei

WORKDIR /root/convert
COPY ./Convert.java Convert.java
COPY ./build.gradle build.gradle
RUN gradle build && \
    cp build/libs/convert.jar /root/convert.jar && \
    rm -rf /root/convert

WORKDIR /root
COPY ./javadoc.conf /root/javadoc.conf
COPY ./compile.sh /root/compile.sh
RUN chmod 777 /root/compile.sh

#CMD /bin/bash
CMD /root/compile.sh
