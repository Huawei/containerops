#!/bin/sh


cat > ../ca/kube-apiserver-admin-ca-config.json <<EOF
{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {
            "kube-apiserver-admin": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            }
        }
    }
}
EOF

cat > ../ca/kube-apiserver-admin-ca-csr.json <<EOF
{
  "CN": "kubernetes-admin",
  "hosts": [
      "10.138.232.140",
      "10.138.48.164",
      "10.138.232.252",
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
      "O": "system:masters",
      "OU": "cloudnative"
    }
  ]
}
EOF

cfssl gencert \
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kube-apiserver-admin-ca-config.json \
    -profile=kube-apiserver-admin ../ca/kube-apiserver-admin-ca-csr.json \
    | cfssljson -bare ../ca/kube-apiserver-admin-ca

