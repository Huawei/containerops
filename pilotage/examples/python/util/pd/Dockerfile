FROM golang:1.7

MAINTAINER xiechuan  <xiechuanj@gmail.com>


ENV HOME=/root
ENV TRAVIS_COVERAGE=1 \
    GOPATH=$HOME/gopath \
    PATH=$HOME/gopath/bin:$PATH \
    TRAVIS_BUILD_DIR=$HOME/gopath/src/github.com/pingcap/pd

RUN go get github.com/mattn/goveralls

COPY *.py ./
CMD python pingcappdtest.py
