FROM docker.io/phusion/baseimage:0.9.21
MAINTAINER dian.li <lidian@huawei.com>
RUN apt-get update && apt-get install -y tar git golang
WORKDIR /var/opt/gopath/
ENV GOPATH /var/opt/gopath
ENV PATH $PATH:$GOROOT/bin:$GOPATH:/bin
#RUN git clone https://github.com/Huawei/containerops.git
COPY component-auto-flow.sh component-auto-flow.sh
COPY main.go main.go
<<<<<<< HEAD
RUN go build main.go
CMD main
=======

RUN go build main.go
CMD flow
>>>>>>> 792ad01fdc25eeba1f89f3099ef992ac7dce7b5f
#CMD ./component-auto-flow.sh


