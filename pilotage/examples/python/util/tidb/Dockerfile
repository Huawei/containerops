FROM golang:1.7

MAINTAINER xiechuan  <xiechuanj@gmail.com>

ENV HOME=/root
ENV GOPATH=$HOME/gopath \
    PATH=$HOME/gopath/bin:$PATH \
    TRAVIS_BUILD_DIR=$HOME/gopath/src/github.com/pingcap/tidb 
#    CURDIR=$HOME/gopath/src/github.com/pingcap/tidb    

#RUN git clone https://github.com/pingcap/tidb.git /root/gopath/src/github.com/pingcap/tidb

RUN apt-get update && apt-get install -y vim
RUN mkdir -p /root/gopath/src/github.com/pingcap


WORKDIR /root/gopath/src/github.com/pingcap

COPY * ./

CMD  python pingcaptidbtest.py
