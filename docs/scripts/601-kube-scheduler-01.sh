#!/bin/sh


export KUBE_APISERVER_PROXY="https://10.138.232.140:19999"

../bin/kubernetes/client/kubectl config set-cluster kubernetes \
--certificate-authority="../ca/cluster-root-ca.pem" \
--client-certificate="../ca/kubernetes-rbac-kube-scheduler-01-ca.pem" \
--client-key="../ca/kubernetes-rbac-kube-scheduler-01-ca-key.pem" \
--server="${KUBE_APISERVER_PROXY}" \
--embed-certs=true \
--kubeconfig=./scheduler01.conf


# set-cluster
../bin/kubernetes/client/kubectl config set-cluster kubernetes \
--certificate-authority="../ca/cluster-root-ca.pem" \
--embed-certs=true \
--server=${KUBE_APISERVER_PROXY} \
--kubeconfig=./scheduler01.conf

# set-credentials
../bin/kubernetes/client/kubectl config set-credentials system:kube-scheduler \
--client-certificate="../ca/kubernetes-rbac-kube-scheduler-01-ca.pem" \
--client-key="../ca/kubernetes-rbac-kube-scheduler-01-ca-key.pem" \
--embed-certs=true \
--kubeconfig=./scheduler01.conf

# set-context
../bin/kubernetes/client/kubectl config set-context system:kube-scheduler@kubernetes \
--cluster=kubernetes \
--user=system:kube-scheduler \
--kubeconfig=./scheduler01.conf


# set default context
../bin/kubernetes/client/kubectl \
config use-context system:kube-scheduler@kubernetes \
--kubeconfig=./scheduler01.conf

../bin/kubernetes/server/kube-scheduler \
--logtostderr=true \
--v=0 \
--master=${kube_api_server_proxy} \
--kubeconfig=./scheduler01.conf \
--leader-elect=true
