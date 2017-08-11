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
export ETCD_NODES=${ETCD_NODE_NAME_00}=https://${ETCD_NODE_IP_00}:2380,${ETCD_NODE_NAME_01}=https://${ETCD_NODE_IP_01}:2380,${ETCD_NODE_NAME_02}=https://${ETCD_NODE_IP_02}:2380
#---------------------------------------------------------------------------------
./bin/etcd/etcd \
  --name=${ETCD_NODE_NAME} \
  --cert-file=${ETCD_SERVER_CA_PREFIX}${NOW_ETCD_NODE_SUFFIX}.pem \
  --key-file=${ETCD_SERVER_CA_PREFIX}${NOW_ETCD_NODE_SUFFIX}-key.pem \
  --peer-cert-file=${ETCD_PEER_CA_PREFIX}${NOW_ETCD_NODE_SUFFIX}.pem \
  --peer-key-file=${ETCD_PEER_CA_PREFIX}${NOW_ETCD_NODE_SUFFIX}-key.pem \
  --trusted-ca-file=${ROOT_CA_PREFIX}.pem \
  --peer-trusted-ca-file=${ROOT_CA_PREFIX}.pem \
  --initial-advertise-peer-urls=https://${ETCD_NODE_IP}:2380 \
  --listen-peer-urls=https://${ETCD_NODE_IP}:2380 \
  --listen-client-urls=https://${ETCD_NODE_IP}:2379,http://127.0.0.1:2379 \
  --advertise-client-urls=https://${ETCD_NODE_IP}:2379 \
  --initial-cluster-token=etcd-cluster-${ETCD_TOKEN} \
  --initial-cluster=${ETCD_NODES} \
  --initial-cluster-state=new \
  --data-dir=./data/etcd-${ETCD_TOKEN}-${NOW_ETCD_NODE_SUFFIX} \
  --enable-v2

