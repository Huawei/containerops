#!/bin/sh


cat > ../ca/kubernetes-rbac-kube-scheduler-00-ca.json <<EOF
{
    "CN": "system:kube-scheduler",
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
          "O": "system:kube-scheduler",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kube-scheduler ../ca/kubernetes-rbac-kube-scheduler-00-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-scheduler-00-ca


#---
cat > ../ca/kubernetes-rbac-kube-scheduler-01-ca.json <<EOF
{
    "CN": "system:kube-scheduler",
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
          "O": "system:kube-scheduler",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kube-scheduler ../ca/kubernetes-rbac-kube-scheduler-01-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-scheduler-01-ca
#---
cat > ../ca/kubernetes-rbac-kube-scheduler-02-ca.json <<EOF
{
    "CN": "system:kube-scheduler",
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
          "O": "system:kube-scheduler",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kube-scheduler ../ca/kubernetes-rbac-kube-scheduler-02-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-scheduler-02-ca

