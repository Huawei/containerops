FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER dian.li <lidian@huawei.com>
RUN apt-get update && apt-get install -y tar git golang
WORKDIR /var/opt/gopath/
ENV GOPATH /var/opt/gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH:/bin
#RUN git clone https://github.com/Huawei/containerops.git
COPY component-auto-flow.sh component-auto-flow.sh
COPY main.go main.go
RUN go build main.go
CMD main
#CMD ./component-auto-flow.sh


