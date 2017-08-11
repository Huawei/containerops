# 当前部署的所有机器 IP
export NODE_IP_00=10.138.48.164
export NODE_IP_01=10.138.232.252
export NODE_IP_02=10.138.24.24
#---------------------------------------------------------------------------------
# 当前部署的机器 PREFIX SUFFIX
export ROOT_CA_PREFIX=./root-ca/root-ca
#---------------------------------------------------------------------------------

# TLS Bootstrapping 使用的 Token，可以使用命令 head -c 16 /dev/urandom | od -An -t x | tr -d ' ' 生成
export BOOTSTRAP_TOKEN="dec0ac166ff2dbf8eab068ca47decaa4"

export INTERNAL_IP=${NODE_IP_00}

# 最好使用 主机未用的网段 来定义服务网段和 Pod 网段

# 服务网段 (Service CIDR），部署前路由不可达，部署后集群内使用IP:Port可达
SERVICE_CIDR="10.254.0.0/16"

# POD 网段 (Cluster CIDR），部署前路由不可达，**部署后**路由可达(flanneld保证)
CLUSTER_CIDR="172.30.0.0/16"

# 服务端口范围 (NodePort Range)
export NODE_PORT_RANGE="8400-9000"

# etcd 集群服务地址列表
export ETCD_ENDPOINTS="https://${NODE_IP_00}:2379,https://${NODE_IP_01}:2379,https://${NODE_IP_02}:2379"

# flanneld 网络配置前缀
export FLANNEL_ETCD_PREFIX="/kubernetes/network"

# kubernetes 服务 IP (一般是 SERVICE_CIDR 中第一个IP)
export CLUSTER_KUBERNETES_SVC_IP="10.254.0.1"

# 集群 DNS 服务 IP (从 SERVICE_CIDR 中预分配)
export CLUSTER_DNS_SVC_IP="10.254.0.2"

# 集群 DNS 域名
export CLUSTER_DNS_DOMAIN="cluster.local."

export kube_bootstrap_tokens_filename="k8s-bootstrap-token"


#cat > ./start.kube.sh <<EOF

../bin/kubernetes/server/kube-apiserver \
--apiserver-count=3 \
--insecure-bind-address=${NODE_IP_00} \
--insecure-port=18080 \
--advertise-address=${INTERNAL_IP} \
--etcd-servers=${ETCD_ENDPOINTS} \
--etcd-cafile=./root-ca/root-ca.pem \
--etcd-certfile=./etcd-ca/etcd-client.pem \
--etcd-keyfile=./etcd-ca/etcd-client-key.pem \
--storage-backend=etcd3 \
--experimental-bootstrap-token-auth=true \
--token-auth-file=./${kube_bootstrap_tokens_filename} \
--authorization-mode=RBAC \
--kubelet-https=true \
--service-cluster-ip-range=${SERVICE_CIDR} \
--service-node-port-range=${NODE_PORT_RANGE} \
--tls-cert-file=./kubernetes-ca/kubernetes-api-server.pem \
--tls-private-key-file=./kubernetes-ca/kubernetes-api-server-key.pem \
--client-ca-file=./root-ca/root-ca.pem \
--service-account-key-file=./root-ca/root-ca-key.pem \
--allow-privileged=true \
--enable-swagger-ui=true \
--admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,PersistentVolumeLabel,DefaultStorageClass,ResourceQuota,DefaultTolerationSeconds \
--audit-log-maxage=30 \
--audit-log-maxbackup=3 \
--audit-log-maxsize=100 \
--audit-log-path=/var/log/kubernetes/audit.log \
--v=0 

#EOF

#下面这两个没明白关键点在哪里
#--experimental-bootstrap-token-auth=true \
#--token-auth-file=/etc/kubernetes/token.csv \

