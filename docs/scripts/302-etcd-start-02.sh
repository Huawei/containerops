#!/bin/sh

../bin/etcd/etcd \
  --name=etcd-02 \
  --cert-file=../ca/etcd-server.pem \
  --key-file=../ca/etcd-server-key.pem \
  --trusted-ca-file=../ca/cluster-root-ca.pem \
  --peer-cert-file=../ca/etcd-peer.pem \
  --peer-key-file=../ca/etcd-peer-key.pem \
  --peer-trusted-ca-file=../ca/cluster-root-ca.pem \
  --initial-advertise-peer-urls=https://10.138.24.24:2380 \
  --listen-peer-urls=https://10.138.24.24:2380 \
  --listen-client-urls=https://10.138.24.24:2379 \
  --advertise-client-urls=https://10.138.24.24:2379 \
  --initial-cluster-token=etcd-cluster-dec0ac166ff2dbf8eab068ca47decaa4 \
  --initial-cluster=etcd-00=https://10.138.48.164:2380,etcd-01=https://10.138.232.252:2380,etcd-02=https://10.138.24.24:2380 \
  --initial-cluster-state=new \
  --data-dir=/var/lib/etcd/etcd-dec0ac166ff2dbf8eab068ca47decaa4-etcd-02 \
  --enable-v2
