#!/bin/sh


cat > ../ca/kubernetes-rbac-ca-config.json <<EOF
{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {
            "rbac": {
                "expiry": "43800h",
                "usages": [
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

cat > ../ca/kubernetes-rbac-ca.json <<EOF
{
    "CN": "rbac",
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
    -ca=../ca/kubernetes-root-ca.pem \
    -ca-key=../ca/kubernetes-root-ca-key.pem \
    -config=../ca/kubernetes-rbac-ca-config.json \
    -profile=rbac ../ca/kubernetes-rbac-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-ca

cat > ../ca/kubernetes-rbac-core-component-roles-ca-config.json <<EOF
{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {
            "kube-controller-manager": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            },
            "kube-scheduler": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            },
            "kubelet-node": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            },
            "kube-proxy": {
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
