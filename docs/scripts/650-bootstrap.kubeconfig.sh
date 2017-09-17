#!/bin/sh

export KUBE_APISERVER_PROXY="https://10.138.232.140:19999"
export BOOTSTRAP_TOKEN="dec0ac166ff2dbf8eab068ca47decaa4"


../bin/kubernetes/client/kubectl config set-cluster kubernetes \
--certificate-authority="../ca/cluster-root-ca.pem" \
--server="${KUBE_APISERVER_PROXY}" \
--embed-certs=true \
--kubeconfig=./bootstrap.kubeconfig

# 设置客户端认证参数
../bin/kubernetes/client/kubectl config set-credentials kubelet-bootstrap \
  --token=${BOOTSTRAP_TOKEN} \
  --kubeconfig=./bootstrap.kubeconfig

# 设置上下文参数
../bin/kubernetes/client/kubectl config set-context default \
  --cluster=kubernetes \
  --user=kubelet-bootstrap \
  --kubeconfig=./bootstrap.kubeconfig

# 设置默认上下文
../bin/kubernetes/client/kubectl config use-context default --kubeconfig=./bootstrap.kubeconfig
