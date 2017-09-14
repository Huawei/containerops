#!/bin/sh


#---rbac-kube-controller-manager-00
cat > ../ca/kubernetes-rbac-kube-controller-manager-00-ca.json <<EOF
{
    "CN": "system:kube-controller-manager",
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
          "O": "system:kube-controller-manager",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kube-controller-manager ../ca/kubernetes-rbac-kube-controller-manager-00-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-controller-manager-00-ca

#---rbac-kube-controller-manager-01
cat > ../ca/kubernetes-rbac-kube-controller-manager-01-ca.json <<EOF
{
    "CN": "system:kube-controller-manager",
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
          "O": "system:kube-controller-manager",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kube-controller-manager ../ca/kubernetes-rbac-kube-controller-manager-01-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-controller-manager-01-ca

#---rbac-kube-controller-manager-02
cat > ../ca/kubernetes-rbac-kube-controller-manager-02-ca.json <<EOF
{
    "CN": "system:kube-controller-manager",
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
          "O": "system:kube-controller-manager",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/kubernetes-rbac-ca.pem \
    -ca-key=../ca/kubernetes-rbac-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kube-controller-manager ../ca/kubernetes-rbac-kube-controller-manager-02-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-controller-manager-02-ca

