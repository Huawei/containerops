FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER dian.li <lidian@huawei.com>
RUN apt-get update && apt-get install -y tar git golang
WORKDIR /var/opt/gopath/src
ENV GOPATH /var/opt/gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH:/bin
RUN git clone https://github.com/Huawei/containerops.git
RUN cd containerops/component/
ADD  main.go main.go
ADD ext ext
ADD component_dir component_dir
COPY  dir_reader.sh dir_reader.sh
RUN ./component-auto-ci.sh


