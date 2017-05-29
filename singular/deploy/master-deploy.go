/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package deploy

import (
	cmd "github.com/Huawei/containerops/singular/cmd"
	"github.com/Huawei/containerops/singular/init_config"
)

func Deploymaster(list map[string]string, ip string) {

	opsfirewalld()

	echoconfig(list)
	opsetcd()
	opskube_apiserver()
	opskube_controller_manager()
	opskube_scheduler()
	opsflanneld()
	// #refresh service
	cmd.Reload()

	// #for SERVICES in etcd kube-apiserver kube-controller-manager kube-scheduler ; do
	// 	#systemctl restart $SERVICES
	// 	#systemctl enable $SERVICES
	// 	#systemctl status $SERVICES
	// #done
	cmd.RestartSvc([]string{"etcd", "kube-apiserver", "kube-controller-manager", "kube-scheduler"})

	opskubectlconfig()

	// // kubectl get nodes
	cmd.ExecCMDparams("kubectl", []string{"config", "cluster-info", "dump"})

}

func opsfirewalld() {

	// firewalld
	cmd.ExecCMDparams("systemctl", []string{"disable", "firewalld"})
	cmd.ExecCMDparams("systemctl", []string{"stop", "firewalld"})
}
func echoconfig(list map[string]string) {

	for k, v := range list {
		cmd.ExecShCommandEcho(k, v)
	}

}
func opsetcd() {

	//#etc create dir
	cmd.ExecCMDparams("mkdir", []string{"/etc/kubernetes/"})
	cmd.ExecCPparams("/tmp/config/config", "/etc/kubernetes/config") //
	// #etc service
	cmd.ExecCMDparams("mkdir", []string{"/etc/etcd/"})
	cmd.ExecCPparams("/tmp/etcd.conf", "/etc/etcd/etcd.conf")

	cmd.ExecCMDparams("mkdir", []string{"/usr/lib/systemd/system/"})
	cmd.ExecCPparams("/tmp/config/etcd.service", "/usr/lib/systemd/system/etcd.service") //
	cmd.ExecCPparams("/tmp/etcd", "/usr/bin/etcd")
	cmd.ExecCPparams("/tmp/etcdctl", "/usr/bin/etcdctl")
	cmd.ExecCMDparams("mkdir", []string{"/var/lib/etcd/"})
	// #/var/lib/etcd/default.etcd auto?

	cmd.ExecCMDparams("systemctl", []string{"daemon-reload"})
	cmd.ServiceStart("etcd")
	cmd.ServiceIsEnabled("etcd")
	cmd.ServiceExists("etcd")
	cmd.ExecCMDparams("etcdctl", []string{"mkdir", init_config.EtcdNet}) //"/kube-centos/network/config nend update config
	cmd.ExecCMDparams("etcdctl", []string{"mk", init_config.EtcdNet, "{\"Network\":\"172.40.0.0/16\"\\,\"SubnetLen\":24\\,\"Backend\":{\"Type\":\"vxlan\"}}"})

	//  sudo mkdir -p /var/lib/etcd  # 必须先创建工作目录
	//  cat > etcd.service <<EOF
	// [Unit]
	// Description=Etcd Server
	// After=network.target
	// After=network-online.target
	// Wants=network-online.target
	// Documentation=https://github.com/coreos

	// [Service]
	// Type=notify
	// WorkingDirectory=/var/lib/etcd/
	// ExecStart=/root/local/bin/etcd \\
	//   --name=${NODE_NAME} \\
	//   --cert-file=/etc/etcd/ssl/etcd.pem \\
	//   --key-file=/etc/etcd/ssl/etcd-key.pem \\
	//   --peer-cert-file=/etc/etcd/ssl/etcd.pem \\
	//   --peer-key-file=/etc/etcd/ssl/etcd-key.pem \\
	//   --trusted-ca-file=/etc/kubernetes/ssl/ca.pem \\
	//   --peer-trusted-ca-file=/etc/kubernetes/ssl/ca.pem \\
	//   --initial-advertise-peer-urls=https://${NODE_IP}:2380 \\
	//   --listen-peer-urls=https://${NODE_IP}:2380 \\
	//   --listen-client-urls=https://${NODE_IP}:2379,http://127.0.0.1:2379 \\
	//   --advertise-client-urls=https://${NODE_IP}:2379 \\
	//   --initial-cluster-token=etcd-cluster-0 \\
	//   --initial-cluster=${ETCD_NODES} \\
	//   --initial-cluster-state=new \\
	//   --data-dir=/var/lib/etcd
	// Restart=on-failure
	// RestartSec=5
	// LimitNOFILE=65536

	// [Install]
	// WantedBy=multi-user.target
	// EOF

	// check
	// 	 for ip in ${NODE_IPS}; do
	//   ETCDCTL_API=3 /root/local/bin/etcdctl \
	//   --endpoints=https://${ip}:2379  \
	//   --cacert=/etc/kubernetes/ssl/ca.pem \
	//   --cert=/etc/etcd/ssl/etcd.pem \
	//   --key=/etc/etcd/ssl/etcd-key.pem \
	//   endpoint health; done
}

func opskube_apiserver() {
	//#kube-apiserver
	cmd.ExecCPparams("/tmp/kube-apiserver", "/usr/bin/kube-apiserver")

	cmd.ExecCPparams("/tmp/config/kube-apiserver.service", "/usr/lib/systemd/system/kube-apiserver.service")
	cmd.ExecCPparams("/tmp/config/apiserver", "/etc/kubernetes/apiserver")

	// 	cat  > kube-apiserver.service <<EOF
	// [Unit]
	// Description=Kubernetes API Server
	// Documentation=https://github.com/GoogleCloudPlatform/kubernetes
	// After=network.target

	// [Service]
	// ExecStart=/root/local/bin/kube-apiserver \\
	//   --admission-control=NamespaceLifecycle,LimitRanger,ServiceAccount,DefaultStorageClass,ResourceQuota \\
	//   --advertise-address=${MASTER_IP} \\
	//   --bind-address=${MASTER_IP} \\
	//   --insecure-bind-address=${MASTER_IP} \\
	//   --authorization-mode=RBAC \\
	//   --runtime-config=rbac.authorization.k8s.io/v1alpha1 \\
	//   --kubelet-https=true \\
	//   --experimental-bootstrap-token-auth \\
	//   --token-auth-file=/etc/kubernetes/token.csv \\
	//   --service-cluster-ip-range=${SERVICE_CIDR} \\
	//   --service-node-port-range=${NODE_PORT_RANGE} \\
	//   --tls-cert-file=/etc/kubernetes/ssl/kubernetes.pem \\
	//   --tls-private-key-file=/etc/kubernetes/ssl/kubernetes-key.pem \\
	//   --client-ca-file=/etc/kubernetes/ssl/ca.pem \\
	//   --service-account-key-file=/etc/kubernetes/ssl/ca-key.pem \\
	//   --etcd-cafile=/etc/kubernetes/ssl/ca.pem \\
	//   --etcd-certfile=/etc/kubernetes/ssl/kubernetes.pem \\
	//   --etcd-keyfile=/etc/kubernetes/ssl/kubernetes-key.pem \\
	//   --etcd-servers=${ETCD_ENDPOINTS} \\
	//   --enable-swagger-ui=true \\
	//   --allow-privileged=true \\
	//   --apiserver-count=3 \\
	//   --audit-log-maxage=30 \\
	//   --audit-log-maxbackup=3 \\
	//   --audit-log-maxsize=100 \\
	//   --audit-log-path=/var/lib/audit.log \\
	//   --event-ttl=1h \\
	//   --v=2
	// Restart=on-failure
	// RestartSec=5
	// Type=notify
	// LimitNOFILE=65536

	// [Install]
	// WantedBy=multi-user.target
	// EOF

}
func opskube_controller_manager() {
	//#kube-controller-manager
	cmd.ExecCPparams("/tmp/kube-controller-manager", "/usr/bin/kube-controller-manager")

	cmd.ExecCPparams("/tmp/config/kube-controller-manager.service", "/usr/lib/systemd/system/kube-controller-manager.service")
	cmd.ExecCPparams("/tmp/config/controller-manager", "/etc/kubernetes/controller-manager")

	//cat > kube-controller-manager.service<<EOF
	// [Unit]
	// Description=Kubernetes Controller Manager
	// Documentation=https://github.com/GoogleCloudPlatform/kubernetes

	// [Service]
	// ExecStart=/root/local/bin/kube-controller-manager \\
	//   --address=127.0.0.1 \\
	//   --master=http://${MASTER_IP}:8080 \\
	//   --allocate-node-cidrs=true \\
	//   --service-cluster-ip-range=${SERVICE_CIDR} \\
	//   --cluster-cidr=${CLUSTER_CIDR} \\
	//   --cluster-name=kubernetes \\
	//   --cluster-signing-cert-file=/etc/kubernetes/ssl/ca.pem \\
	//   --cluster-signing-key-file=/etc/kubernetes/ssl/ca-key.pem \\
	//   --service-account-private-key-file=/etc/kubernetes/ssl/ca-key.pem \\
	//   --root-ca-file=/etc/kubernetes/ssl/ca.pem \\
	//   --leader-elect=true \\
	//   --v=2
	// Restart=on-failure
	// RestartSec=5

	// [Install]
	// WantedBy=multi-user.target
	// EOF

}
func opskube_scheduler() {

	//#kube-scheduler
	cmd.ExecCPparams("/tmp/kube-scheduler", "/usr/bin/kube-scheduler")

	cmd.ExecCPparams("/tmp/config/kube-scheduler.service", "/usr/lib/systemd/system/kube-scheduler.service")
	cmd.ExecCPparams("/tmp/config/scheduler", "/etc/kubernetes/scheduler")

	// 	cat > kube-scheduler.service <<EOF
	// [Unit]
	// Description=Kubernetes Scheduler
	// Documentation=https://github.com/GoogleCloudPlatform/kubernetes

	// [Service]
	// ExecStart=/root/local/bin/kube-scheduler \\
	//   --address=127.0.0.1 \\
	//   --master=http://${MASTER_IP}:8080 \\
	//   --leader-elect=true \\
	//   --v=2
	// Restart=on-failure
	// RestartSec=5

	// [Install]
	// WantedBy=multi-user.target
	// EOF

}
func opsflanneld() {
	// #flanneld reserved
	cmd.ExecCPparams("/tmp/flanneld", "/usr/bin/flanneld")

	cmd.ExecCMDparams("mkdir", []string{"-p", "/usr/libexec/flannel/"})
	cmd.ExecCMDparams("mkdir", []string{"/tmp/etcd"})
	// check work
	cmd.ExecCPparams("/tmp/flannel/usr/libexec/flannel/mk-docker-opts.sh", "/usr/libexec/flannel/mk-docker-opts.sh")
	cmd.ExecCMDparams("mkdir", []string{"-p", "/etc/sysconfig/"})

	cmd.ExecCPparams("/tmp/config/flanneld", "/etc/sysconfig/flanneld")

	cmd.ExecCPparams("/tmp/config/flanneld.service", "/usr/lib/systemd/system/flanneld.service")

	// flanneld.service
	//  cat > flanneld.service << EOF
	// [Unit]
	// Description=Flanneld overlay address etcd agent
	// After=network.target
	// After=network-online.target
	// Wants=network-online.target
	// After=etcd.service
	// Before=docker.service

	// [Service]
	// Type=notify
	// ExecStart=/root/local/bin/flanneld \\
	//   -etcd-cafile=/etc/kubernetes/ssl/ca.pem \\
	//   -etcd-certfile=/etc/flanneld/ssl/flanneld.pem \\
	//   -etcd-keyfile=/etc/flanneld/ssl/flanneld-key.pem \\
	//   -etcd-endpoints=${ETCD_ENDPOINTS} \\
	//   -etcd-prefix=${FLANNEL_ETCD_PREFIX}
	// ExecStartPost=/root/local/bin/mk-docker-opts.sh -k DOCKER_NETWORK_OPTIONS -d /run/flannel/docker
	// Restart=on-failure

	// [Install]
	// WantedBy=multi-user.target
	// RequiredBy=docker.service
	// EOF

	// 	 /root/local/bin/etcdctl \
	//   --endpoints=${ETCD_ENDPOINTS} \
	//   --ca-file=/etc/kubernetes/ssl/ca.pem \
	//   --cert-file=/etc/flanneld/ssl/flanneld.pem \
	//   --key-file=/etc/flanneld/ssl/flanneld-key.pem \
	//   set ${FLANNEL_ETCD_PREFIX}/config '{"Network":"'${CLUSTER_CIDR}'", "SubnetLen": 24, "Backend": {"Type": "vxlan"}}'

}

func opskubectlconfig() {
	// #kubectl config
	cmd.ExecCPparams("/tmp/kubectl", "/usr/bin/kubectl")

	cmd.ExecCMDparams("kubectl", []string{"config", "set-cluster", "default-cluster", "--server=http://centos-master:8080"})
	cmd.ExecCMDparams("kubectl", []string{"config", "set-cluster", "default-cluster", "--cluster=default-cluster", "--user=default-admin"})
	cmd.ExecCMDparams("kubectl", []string{"config", "use-context", "use-context"})
	cmd.ExecCMDparams("kubectl", []string{"config", "get", "nodes"})

}
