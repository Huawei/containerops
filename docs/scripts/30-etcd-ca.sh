#!/bin/sh

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

cat > ../ca/etcd-server-00-ca.json <<EOF
{
    "CN": "etcd-server-00",
    "hosts": ["10.138.48.164"],
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

cat > ../ca/etcd-server-01-ca.json <<EOF
{
    "CN": "etcd-server-01",
    "hosts": ["10.138.232.252"],
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

cat > ../ca/etcd-server-02-ca.json <<EOF
{
    "CN": "etcd-server-02",
    "hosts": ["10.138.24.24"],
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
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
    -config=../ca/etcd-ca-config.json \
    -profile=server ../ca/etcd-server-00-ca.json \
    | cfssljson -bare ../ca/etcd-server-00-ca

cfssl gencert \
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
    -config=../ca/etcd-ca-config.json \
    -profile=server ../ca/etcd-server-01-ca.json \
    | cfssljson -bare ../ca/etcd-server-01-ca

cfssl gencert \
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
    -config=../ca/etcd-ca-config.json \
    -profile=server ../ca/etcd-server-02-ca.json \
    | cfssljson -bare ../ca/etcd-server-02-ca

cat > ../ca/etcd-peer-00-ca.json <<EOF
{
    "CN": "etcd-peer-00",
    "hosts": ["10.138.48.164"],
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

cat > ../ca/etcd-peer-01-ca.json <<EOF
{
    "CN": "etcd-peer-01",
    "hosts": ["10.138.232.252"],
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

cat > ../ca/etcd-peer-02-ca.json <<EOF
{
    "CN": "etcd-peer-02",
    "hosts": ["10.138.24.24"],
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
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
    -config=../ca/etcd-ca-config.json \
    -profile=peer ../ca/etcd-peer-00-ca.json \
    | cfssljson -bare ../ca/etcd-peer-00-ca

cfssl gencert \
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
    -config=../ca/etcd-ca-config.json \
    -profile=peer ../ca/etcd-peer-01-ca.json \
    | cfssljson -bare ../ca/etcd-peer-01-ca

cfssl gencert \
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
    -config=../ca/etcd-ca-config.json \
    -profile=peer ../ca/etcd-peer-02-ca.json \
    | cfssljson -bare ../ca/etcd-peer-02-ca


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
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
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
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
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
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
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
    -ca=../ca/etcd-root-ca.pem \
    -ca-key=../ca/etcd-root-ca-key.pem \
    -config=../ca/etcd-ca-config.json \
    -profile=client ../ca/etcd-client-other-general-ca.json \
    | cfssljson -bare ../ca/etcd-client-other-general-ca


