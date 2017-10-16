#!/bin/sh



export KUBE_APISERVER_PROXY="https://10.138.232.140:19999"

../bin/kubernetes/client/kubectl config set-cluster kubernetes \
--certificate-authority="../ca/cluster-root-ca.pem" \
--client-certificate="../ca/kubernetes-rbac-kube-controller-manager-02-ca.pem" \
--client-key="../ca/kubernetes-rbac-kube-controller-manager-02-ca-key.pem" \
--server="${KUBE_APISERVER_PROXY}" \
--embed-certs=true \
--kubeconfig=controller-manager02.conf


# set-cluster
../bin/kubernetes/client/kubectl config set-cluster kubernetes \
--certificate-authority="../ca/cluster-root-ca.pem" \
--embed-certs=true \
--server=${KUBE_APISERVER_PROXY} \
--kubeconfig=controller-manager02.conf

# set-credentials
../bin/kubernetes/client/kubectl config set-credentials system:kube-controller-manager \
--client-certificate="../ca/kubernetes-rbac-kube-controller-manager-02-ca.pem" \
--client-key="../ca/kubernetes-rbac-kube-controller-manager-02-ca-key.pem" \
--embed-certs=true \
--kubeconfig=controller-manager02.conf

# set-context
../bin/kubernetes/client/kubectl config set-context system:kube-controller-manager@kubernetes \
--cluster=kubernetes \
--user=system:kube-controller-manager \
--kubeconfig=controller-manager02.conf


# set default context
../bin/kubernetes/client/kubectl config use-context system:kube-controller-manager@kubernetes --kubeconfig=controller-manager02.conf


../bin/kubernetes/server/kube-controller-manager \
--logtostderr=true \
--v=0 \
--master=${KUBE_APISERVER} \
--kubeconfig=./controller-manager02.conf \
--cluster-name=kubernetes \
--cluster-signing-cert-file=../ca/cluster-root-ca.pem \
--cluster-signing-key-file=../ca/cluster-root-ca-key.pem \
--service-account-private-key-file=../ca/cluster-root-ca-key.pem \
--root-ca-file=../ca/cluster-root-ca.pem \
--insecure-experimental-approve-all-kubelet-csrs-for-group=system:bootstrappers \
--use-service-account-credentials=true \
--service-cluster-ip-range=10.254.0.0/16 \
--cluster-cidr=172.30.0.0/16 \
--allocate-node-cidrs=true \
--leader-elect=true \
--controllers=*,bootstrapsigner,tokencleaner



