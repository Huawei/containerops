#!/bin/sh

#--生成root ca csr 证书--------------------------------------------------------
rm -rf root-ca
mkdir root-ca

cat > ./root-ca/root-ca-csr.json <<EOF
{
  "CN": "etcd-peer-00",
  "hosts": [
    "127.0.0.1",
    "${ETCD_PEER_00_IP}"
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

cfssl gencert -initca ./root-ca/root-ca-csr.json | cfssljson -bare ./root-ca/root-ca
#----------------------------------------------------------
# 要部署ETCD的机器 IP
export ETCD_PEER_00_IP=10.138.48.164
export ETCD_PEER_01_IP=10.138.232.252
export ETCD_PEER_02_IP=10.138.24.24
#----------------------------------------------------------
rm -rf ./etcd-ca
mkdir ./etcd-ca
#----------------------------------------------------------
cat > ./etcd-ca/etcd-peer-00-csr.json <<EOF
{
  "CN": "etcd-peer-00",
  "hosts": [
    "127.0.0.1",
    "${ETCD_PEER_00_IP}"
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

cfssl gencert -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=etcd-peer ./etcd-ca/etcd-peer-00-csr.json | cfssljson -bare ./etcd-ca/etcd-peer-00
#----------------------------------------------------------
cat > ./etcd-ca/etcd-peer-01-csr.json <<EOF
{
  "CN": "etcd-peer-01",
  "hosts": [
    "127.0.0.1",
    "${ETCD_PEER_01_IP}"
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

cfssl gencert -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=etcd-peer ./etcd-ca/etcd-peer-01-csr.json | cfssljson -bare ./etcd-ca/etcd-peer-01
#----------------------------------------------------------
cat > ./etcd-ca/etcd-peer-02-csr.json <<EOF
{
  "CN": "etcd-peer-02",
  "hosts": [
    "127.0.0.1",
    "${ETCD_PEER_02_IP}"
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

cfssl gencert -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=etcd-peer ./etcd-ca/etcd-peer-02-csr.json | cfssljson -bare ./etcd-ca/etcd-peer-02

#----------------------------------------------------------

cat > ./etcd-ca/etcd-server00.json <<EOF
{
    "CN": "etcd-server00",
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

cfssl gencert \
    -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=etcd-server ./etcd-ca/etcd-server00.json \
    | cfssljson -bare ./etcd-ca/etcd-server00

#----------------------------------------------------------
cat > ./etcd-ca/etcd-server01.json <<EOF
{
    "CN": "etcd-server01",
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

cfssl gencert \
    -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=etcd-server ./etcd-ca/etcd-server01.json \
    | cfssljson -bare ./etcd-ca/etcd-server01

#----------------------------------------------------------
cat > ./etcd-ca/etcd-server02.json <<EOF
{
    "CN": "etcd-server02",
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
    -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=etcd-server ./etcd-ca/etcd-server02.json \
    | cfssljson -bare ./etcd-ca/etcd-server02

#----------------------------------------------------------

cat > ./etcd-ca/etcd-client.json <<EOF
{
    "CN": "etcd-client",
    "hosts": [""],
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
    -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=etcd-client ./etcd-ca/etcd-client.json \
    | cfssljson -bare ./etcd-ca/etcd-client

#----------------------------------------------------------
rm -rf ./flannel-ca
mkdir ./flannel-ca

cat > ./flannel-ca/flannel.json <<EOF
{
  "CN": "flannel",
  "hosts": [],
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
    -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=flannel ./flannel-ca/flannel.json \
    | cfssljson -bare ./flannel-ca/flannel
#----------------------------------------------------------



#--生成访问k8s集群的master证书--------------------------------------------------------
rm -rf ./kubernetes-ca
mkdir ./kubernetes-ca
#----------------------------------------------------------
cat > ./kubernetes-ca/kubernetes-admin.json <<EOF
{
  "CN": "admin",
  "hosts": [],
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
    -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=kubernetes-admin ./kubernetes-ca/kubernetes-admin.json \
    | cfssljson -bare ./kubernetes-ca/kubernetes-admin
#----------------------------------------------------------
cat > ./kubernetes-ca/kubernetes-api-server.json <<EOF
{
  "CN": "kubernetes-api-server",
  "hosts": [
    "127.0.0.1",
    "localhost",
    "10.138.48.164",
    "10.138.232.252",
    "10.138.24.24",
    "138.197.216.239",
    "165.227.28.28",
    "165.227.28.41",
    "kubernetes",
    "kubernetes.default",
    "kubernetes.default.svc",
    "kubernetes.default.svc.cluster",
    "kubernetes.default.svc.cluster.local"
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

cfssl gencert \
    -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=kubernetes-api-server ./kubernetes-ca/kubernetes-api-server.json \
    | cfssljson -bare ./kubernetes-ca/kubernetes-api-server

#----------------------------------------------------------
    cat > ./kubernetes-ca/kubernetes-controller-manager.json <<EOF
    {
      "CN": "kubernetes-controller-manager",
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
          "OU": "System"
        }
      ]
    }
EOF

    cfssl gencert \
        -ca=./root-ca/root-ca.pem \
        -ca-key=./root-ca/root-ca-key.pem \
        -config=./ca-config.json \
        -profile=kubernetes-controller-manager ./kubernetes-ca/kubernetes-controller-manager.json \
        | cfssljson -bare ./kubernetes-ca/kubernetes-controller-manager

#----------------------------------------------------------
cat > ./kubernetes-ca/kubernetes-scheduler.json <<EOF
{
  "CN": "kubernetes-scheduler",
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
      "OU": "System"
    }
  ]
}
EOF

cfssl gencert \
    -ca=./root-ca/root-ca.pem \
    -ca-key=./root-ca/root-ca-key.pem \
    -config=./ca-config.json \
    -profile=kubernetes-scheduler ./kubernetes-ca/kubernetes-scheduler.json \
    | cfssljson -bare ./kubernetes-ca/kubernetes-scheduler
#----------------------------------------------------------
#----------------------------------------------------------
#----------------------------------------------------------
#----------------------------------------------------------
#----------------------------------------------------------
#----------------------------------------------------------
#----------------------------------------------------------

