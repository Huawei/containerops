FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER dian.li <lidian@huawei.com>
RUN apt-get update && apt-get install -y tar git golang
WORKDIR /var/opt/gopath/src
ENV GOPATH /var/opt/gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH:/bin
RUN mkdir -p github.com/Huawei/containerops/component/test/build
RUN git clone https://github.com/Huawei/containerops.git
RUN cd github.com/Huawei/containerops/component/test/build
ADD main.go main.go
ADD module module
COPY component-auto-ci.sh component-auto-ci.sh
CMD ./component-auto-ci.sh


