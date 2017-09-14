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

