#!/bin/sh


cat > ../ca/kubernetes-rbac-kube-proxy-00-ca.json <<EOF
{
    "CN": "system:kube-proxy",
    "hosts": [
        "10.138.48.164"
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
          "O": "system:kube-proxy",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kubelet-node ../ca/kubernetes-rbac-kube-proxy-00-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-proxy-00-ca
#---
cat > ../ca/kubernetes-rbac-kube-proxy-01-ca.json <<EOF
{
    "CN": "system:kube-proxy",
    "hosts": [
        "10.138.232.252"
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
          "O": "system:kube-proxy",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kubelet-node ../ca/kubernetes-rbac-kube-proxy-01-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-proxy-01-ca

#---
cat > ../ca/kubernetes-rbac-kube-proxy-02-ca.json <<EOF
{
    "CN": "system:kube-proxy",
    "hosts": [
        "10.138.24.24"
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
          "O": "system:kube-proxy",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kubelet-node ../ca/kubernetes-rbac-kube-proxy-02-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-proxy-02-ca
#---

