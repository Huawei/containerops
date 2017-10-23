#!/bin/sh


cat > ../ca/kube-apiserver-ca-config.json <<EOF
{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {
            "kube-apiserver": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth"
                ]
            }
        }
    }
}
EOF

cat > ../ca/kube-apiserver-00-ca.json <<EOF
{
    "CN": "kube-apiserver-00",
    "hosts": [
    "10.138.48.164",
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster",
    "kubernetes.default.svc.cluster.local",
    "10.138.232.140",
    "10.254.0.1"
    ],
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

cat > ../ca/kube-apiserver-01-ca.json <<EOF
{
    "CN": "kube-apiserver-01",
    "hosts": [
    "10.138.232.252",
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster",
    "kubernetes.default.svc.cluster.local",
    "10.138.232.140",
    "10.254.0.1"],
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

cat > ../ca/kube-apiserver-02-ca.json <<EOF
{
    "CN": "kube-apiserver-02",
    "hosts": [
    "10.138.24.24",
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster",
    "kubernetes.default.svc.cluster.local",
    "10.138.232.140",
    "10.254.0.1"],
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
    -config=../ca/kube-apiserver-ca-config.json \
    -profile=kube-apiserver ../ca/kube-apiserver-00-ca.json \
    | cfssljson -bare ../ca/kube-apiserver-00-ca

cfssl gencert \
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kube-apiserver-ca-config.json \
    -profile=kube-apiserver ../ca/kube-apiserver-01-ca.json \
    | cfssljson -bare ../ca/kube-apiserver-01-ca

cfssl gencert \
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kube-apiserver-ca-config.json \
    -profile=kube-apiserver ../ca/kube-apiserver-02-ca.json \
    | cfssljson -bare ../ca/kube-apiserver-02-ca


