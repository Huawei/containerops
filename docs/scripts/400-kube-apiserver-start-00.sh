#!/bin/sh


../bin/kubernetes/server/kube-apiserver \
--apiserver-count=3 \
--advertise-address=10.138.48.164 \
--etcd-servers=https://10.138.48.164:2379,https://10.138.232.252:2379,https://10.138.24.24:2379 \
--etcd-cafile=../ca/cluster-root-ca.pem \
--etcd-certfile=../ca/etcd-client-kubernetes-ca.pem \
--etcd-keyfile=../ca/etcd-client-kubernetes-ca-key.pem \
--storage-backend=etcd3 \
--experimental-bootstrap-token-auth=true \
--token-auth-file=./k8s-bootstrap-token \
--authorization-mode=RBAC \
--kubelet-https=true \
--service-cluster-ip-range=10.254.0.0/16  \
--service-node-port-range=8400-9000 \
--tls-cert-file=../ca/kube-apiserver-00-ca.pem \
--tls-private-key-file=../ca/kube-apiserver-00-ca-key.pem \
--client-ca-file=../ca/cluster-root-ca.pem \
--service-account-key-file=../ca/cluster-root-ca-key.pem \
--allow-privileged=true \
--enable-swagger-ui=true \
--admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,ResourceQuota,DefaultTolerationSeconds \
--audit-log-maxage=30 \
--audit-log-maxbackup=3 \
--audit-log-maxsize=100 \
--audit-log-path=/var/log/kubernetes/audit.log \
--v=0

