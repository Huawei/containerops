#!/bin/sh


# kube-apiserver haproxy IP
export KUBE_APISERVER_PROXY="https://10.138.232.140:19999"
export HOST_IP="10.138.232.252"
export BOOTSTRAP_TOKEN="dec0ac166ff2dbf8eab068ca47decaa4"

# 服务网段 (Service CIDR），部署前路由不可达，部署后集群内使用 IP:Port 可达
SERVICE_CIDR="10.254.0.0/16"
# POD 网段 (Cluster CIDR），部署前路由不可达，**部署后**路由可达 (网络插件 保证)
CLUSTER_CIDR="172.30.0.0/16"
# 集群 DNS 服务 IP (从 SERVICE_CIDR 中预分配)
export CLUSTER_DNS_SVC_IP="10.254.0.2"
# 集群 DNS 域名
export CLUSTER_DNS_DOMAIN="cluster.local."

# 配置集群

../bin/kubernetes/client/kubectl config set-cluster kubernetes \
  --certificate-authority=../ca/cluster-root-ca.pem \
  --embed-certs=true \
  --server=${KUBE_APISERVER_PROXY} \
  --kubeconfig=./kube-proxy01.kubeconfig


# 配置客户端认证

../bin/kubernetes/client/kubectl config set-credentials kube-proxy \
  --client-certificate=../ca/kubernetes-rbac-kube-proxy-01-ca.pem \
  --client-key=../ca/kubernetes-rbac-kube-proxy-01-ca-key.pem \
  --embed-certs=true \
  --kubeconfig=./kube-proxy01.kubeconfig


# 配置关联

../bin/kubernetes/client/kubectl config set-context default \
  --cluster=kubernetes \
  --user=kube-proxy \
  --kubeconfig=./kube-proxy01.kubeconfig



# 配置默认关联
../bin/kubernetes/client/kubectl config use-context default --kubeconfig=./kube-proxy01.kubeconfig

../bin/kubernetes/server/kube-proxy \
--bind-address=${HOST_IP} \
--hostname-override=${HOST_IP} \
--cluster-cidr=${SERVICE_CIDR} \
--kubeconfig=./kube-proxy01.kubeconfig \
--logtostderr=true \
--v=0


