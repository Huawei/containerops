#!/bin/sh

cat > ../ca/etcd-client-calico-ca.json <<EOF
{
    "CN": "etcd-client-calico",
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
    -config=../ca/etcd-ca-config.json \
    -profile=client ../ca/etcd-client-calico-ca.json \
    | cfssljson -bare ../ca/etcd-client-calico-ca

cat > ../ca/etcd-client-kubernetes-ca.json <<EOF
{
    "CN": "etcd-client-kubernetes",
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
    -config=../ca/etcd-ca-config.json \
    -profile=client ../ca/etcd-client-kubernetes-ca.json \
    | cfssljson -bare ../ca/etcd-client-kubernetes-ca

cat > ../ca/etcd-client-kubedns-ca.json <<EOF
{
    "CN": "etcd-client-kubedns",
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
    -config=../ca/etcd-ca-config.json \
    -profile=client ../ca/etcd-client-kubedns-ca.json \
    | cfssljson -bare ../ca/etcd-client-kubedns-ca

cat > ../ca/etcd-client-other-general-ca.json <<EOF
{
    "CN": "etcd-client-other-general",
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
    -config=../ca/etcd-ca-config.json \
    -profile=client ../ca/etcd-client-other-general-ca.json \
    | cfssljson -bare ../ca/etcd-client-other-general-ca

