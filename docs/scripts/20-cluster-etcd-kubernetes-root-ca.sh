#!/bin/sh

cat > ../ca/cluster-etcd-kubernetes-root-ca-config.json <<EOF
{
  "signing":{
    "default":{
      "expiry":"43800h"
    },
    "profiles":{
      "etcd-root":{
        "expiry":"43800h",
        "usages":[
          "signing",
          "key encipherment",
          "server auth",
          "client auth"
        ]
      },
      "kubernetes-root":{
        "expiry":"43800h",
        "usages":[
          "signing",
          "key encipherment",
          "server auth",
          "client auth"
        ]
      }
    }
  }
}
EOF

cat > ../ca/etcd-root-ca-csr.json <<EOF
{
  "CN": "etcd-root-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
EOF


cat > ../ca/kubernetes-root-ca-csr.json <<EOF
{
  "CN": "kubernetes-root-ca",
  "key": {
    "algo": "rsa",
    "size": 2048
  },
  "names": [
    {
      "C": "CN",
      "ST": "BeiJing",
      "L": "BeiJing",
      "O": "k8s",
      "OU": "System"
    }
  ]
}
EOF


cfssl gencert \
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/cluster-etcd-kubernetes-root-ca-config.json \
    -profile=etcd-root ../ca/etcd-root-ca-csr.json \
    | cfssljson -bare ../ca/etcd-root-ca

cfssl gencert \
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/cluster-etcd-kubernetes-root-ca-config.json \
    -profile=kubernetes-root ../ca/kubernetes-root-ca-csr.json \
    | cfssljson -bare ../ca/kubernetes-root-ca


