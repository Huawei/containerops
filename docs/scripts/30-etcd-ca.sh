#!/bin/sh

cat > ../ca/cluster-root-ca-csr.json <<EOF
{
  "CN": "cluster-root-ca",
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

cfssl gencert -initca ../ca/cluster-root-ca-csr.json | cfssljson -bare ../ca/cluster-root-ca


cat > ../ca/etcd-ca-config.json <<EOF
{
    "signing": {
        "default": {
            "expiry": "43800h"
        },
        "profiles": {
            "server": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "server auth"
                ]
            },
            "client": {
                "expiry": "43800h",
                "usages": [
                    "signing",
                    "key encipherment",
                    "client auth"
                ]
            },
            "peer": {
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


cat > ../ca/etcd-server-csr.json <<EOF
{
    "CN": "server",
    "hosts": [
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
            "O": "k8s",
            "OU": "cloudnative"
        }
    ]
}
EOF

cfssl gencert \
-ca=../ca/cluster-root-ca.pem \
-ca-key=../ca/cluster-root-ca-key.pem \
-config=../ca/etcd-ca-config.json \
-profile=server ../ca/etcd-server-csr.json | cfssljson -bare ../ca/etcd-server


cat > ../ca/etcd-peer-csr.json <<EOF
{
    "CN": "peer",
    "hosts": [
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
            "O": "k8s",
            "OU": "cloudnative"
        }
    ]
}
EOF



cfssl gencert \
-ca=../ca/cluster-root-ca.pem \
-ca-key=../ca/cluster-root-ca-key.pem \
-config=../ca/etcd-ca-config.json \
-profile=peer ../ca/etcd-peer-csr.json | cfssljson -bare ../ca/etcd-peer


