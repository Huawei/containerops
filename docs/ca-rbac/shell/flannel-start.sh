#!/bin/sh
#---------------------------------------------------------------------------------
export ETCD_TOKEN=dec0ac166ff2dbf8eab068ca47decaa4
#---------------------------------------------------------------------------------
# 当前部署的所有机器 IP
export NODE_IP_00=10.138.48.164
export NODE_IP_01=10.138.232.252
export NODE_IP_02=10.138.24.24
#---------------------------------------------------------------------------------
# 当前部署的机器 PREFIX SUFFIX
export ROOT_CA_PREFIX=./root-ca/root-ca
#---------------------------------------------------------------------------------
# 当前部署ETCD的所有机器 IP
export ETCD_NODE_IP_00=10.138.48.164
export ETCD_NODE_IP_01=10.138.232.252
export ETCD_NODE_IP_02=10.138.24.24
#---------------------------------------------------------------------------------
# 当前部署的所有机器名称(随便定义，只要能区分不同机器即可)
export ETCD_NODE_NAME_00=etcd-host-00
export ETCD_NODE_NAME_01=etcd-host-01
export ETCD_NODE_NAME_02=etcd-host-02
#---------------------------------------------------------------------------------
# 当前ETCD Node有效的Name和IP
export NOW_ETCD_NODE_SUFFIX=00
export ETCD_NODE_NAME=${ETCD_NODE_NAME_00}
export ETCD_NODE_IP=${ETCD_NODE_IP_00}
#---------------------------------------------------------------------------------
# 当前Node有效的ETCD ca证书为止前缀
export ETCD_PEER_CA_PREFIX=./etcd-ca/etcd-peer-
export ETCD_SERVER_CA_PREFIX=./etcd-ca/etcd-server
#---------------------------------------------------------------------------------
# etcd 集群所有机器 IP
export ETCD_NODE_IPS="${ETCD_NODE_IP_00},${ETCD_NODE_IP_01},${ETCD_NODE_IP_02}"
# etcd 集群间通信的IP和端口
export ETCD_ENDPOINTS="https://${ETCD_NODE_IP_00}:2379,https://${ETCD_NODE_IP_01}:2379,https://${ETCD_NODE_IP_02}:2379"
#---------------------------------------------------------------------------------
export FLANNEL_ETCD_PREFIX="/kubernetes/network"
export CLUSTER_CIDR="172.30.0.0/16"

#---------------------------------------------------------------------------------
./bin/etcd/etcdctl \
--cert-file ./etcd-ca/etcd-client.pem \
--key-file ./etcd-ca/etcd-client-key.pem \
--ca-file ./root-ca/root-ca.pem \
--endpoints ${ETCD_ENDPOINTS} \
set ${FLANNEL_ETCD_PREFIX}/config '{"Network":"'${CLUSTER_CIDR}'", "SubnetLen": 24, "Backend": {"Type": "vxlan"}}'
#---------------------------------------------------------------------------------
#---------------------------------------------------------------------------------
./bin/flannel/flanneld \
-etcd-cafile=./root-ca/root-ca.pem \
-etcd-certfile=./flannel-ca/flannel.pem \
-etcd-keyfile=./flannel-ca/flannel-key.pem \
-etcd-endpoints=${ETCD_ENDPOINTS} \
-etcd-prefix=${FLANNEL_ETCD_PREFIX}
#---------------------------------------------------------------------------------
#---------------------------------------------------------------------------------

