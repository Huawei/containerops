#!/bin/sh


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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kube-controller-manager ../ca/kubernetes-rbac-kube-controller-manager-02-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-controller-manager-02-ca

echo --------------------rbac-kube-controller-manager-end--------------------


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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kube-scheduler ../ca/kubernetes-rbac-kube-scheduler-02-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-scheduler-02-ca


echo ---------------------rbac-kube-scheduler-end--------------------

cat > ../ca/kubernetes-rbac-kubelet-node-00-ca.json <<EOF
{
    "CN": "system:node:node00",
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
          "O": "system:nodes",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kubelet-node ../ca/kubernetes-rbac-kubelet-node-00-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kubelet-node-00-ca
#---
cat > ../ca/kubernetes-rbac-kubelet-node-01-ca.json <<EOF
{
    "CN": "system:node:node01",
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
          "O": "system:nodes",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kubelet-node ../ca/kubernetes-rbac-kubelet-node-01-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kubelet-node-01-ca
#---
cat > ../ca/kubernetes-rbac-kubelet-node-02-ca.json <<EOF
{
    "CN": "system:node:node02",
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
          "O": "system:nodes",
          "OU": "System"
        }
    ]
}
EOF

cfssl gencert \
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kubelet-node ../ca/kubernetes-rbac-kubelet-node-02-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kubelet-node-02-ca


echo ---------------------rbac-kubelet-node-end--------------------


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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
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
    -ca=../ca/cluster-root-ca.pem \
    -ca-key=../ca/cluster-root-ca-key.pem \
    -config=../ca/kubernetes-rbac-core-component-roles-ca-config.json \
    -profile=kubelet-node ../ca/kubernetes-rbac-kube-proxy-02-ca.json \
    | cfssljson -bare ../ca/kubernetes-rbac-kube-proxy-02-ca
#---

echo ---------------------rbac-kube-proxy-end--------------------
