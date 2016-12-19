FROM golang
MAINTAINER Marx Guo "go.morefork@gmail.com"
RUN "mkdir" "-p" "/go/src/github.com/kubernetes"
WORKDIR /go/src/github.com/kubernetes
RUN "git" "clone" "https://github.com/kubernetes/kubernetes.git"
ADD ./generateGodoc.go /go/src/
WORKDIR /go/src
CMD go run generateGodoc.go
