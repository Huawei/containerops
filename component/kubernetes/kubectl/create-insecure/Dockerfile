FROM alpine:3.4
MAINTAINER wang qilin <qilin.wang@huawei.com>

ENV KUBECTL_VERSION v1.7.4
ENV EDITOR vim

RUN apk add --no-cache --update ca-certificates wget curl bash vim \
  && wget -qO /usr/local/bin/kubectl "https://storage.googleapis.com/kubernetes-release/release/${KUBECTL_VERSION}/bin/linux/amd64/kubectl" \
  && chmod +x /usr/local/bin/kubectl \
  && apk del --purge wget \
  && rm /var/cache/apk/*

RUN apk add --no-cache \
            --repository http://dl-3.alpinelinux.org/alpine/edge/community/ \
            emacs

WORKDIR /root
COPY run.sh /root/run.sh
RUN chmod 777 /root/run.sh

CMD bash run.sh

