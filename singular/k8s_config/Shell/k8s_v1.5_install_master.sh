

#config your IP
#echo "192.168.121.9	centos-master 192.168.121.65	centos-minion" >> /etc/hosts
echo -e "192.168.121.9 centos-master\n192.168.121.65 centos-minion" >> /etc/hosts
#firewalld
systemctl disable firewalld
systemctl stop firewalld

#change to binary way
yum install docker

#downlload  binary & config to temp
/tmp/etcd
/tmp/k8s_binary
scp -r /Users/dean/Documents/gopath/src/github.com/Huawei/containerops_he3io/singular/k8s_binary/ root@192.168.60.145:/tmp/k8s_binary/
scp -r /Users/dean/Documents/gopath/src/github.com/Huawei/containerops_he3io/singular/k8s_CentOS_Manaul_Install root@192.168.60.145:/tmp/
scp -r /Users/dean/Documents/gopath/src/github.com/Huawei/containerops_he3io/singular/k8s_CentOS_Manaul_Install/config1/ root@192.168.60.145:/tmp/
scp -r /Users/dean/Documents/gopath/src/github.com/Huawei/containerops_he3io/singular/k8s_CentOS_Manaul_Install/. root@192.168.60.141:/tmp/

#erc create dir  
mkdir /etc/kubernetes/
cp /tmp/config/config /etc/kubernetes/config

#etc service
	cp /tmp/etcdetcd.conf /etc/etcd/etcd.conf
	cp /tmp/etcd/etcd.service /usr/lib/systemd/system/etcd.service
	cp /tmp/etcd/usr/bin/etcd  /usr/bin/etcd
	cp /tmp/etcd/usr/bin/etcdctl  /usr/bin/etcdctl
	mkdir /var/lib/etcd/
	#/var/lib/etcd/default.etcd auto?
 
	systemctl daemon-reload
	systemctl start etcd
	systemctl enable etcd
	etcdctl mkdir /kube-centos/network
	etcdctl mk /kube-centos/network/config "{ \"Network\": \"172.40.0.0/16\", \"SubnetLen\": 24, \"Backend\": { \"Type\": \"vxlan\" } }"

#kube-apiserver
	cp /tmp/k8s_binary/kube-apiserver /usr/bin/kube-apiserver
	cp /tmp/config/kube-apiserver.service /usr/lib/systemd/system/kube-apiserver.service 
	cp /tmp/config/apiserver /etc/kubernetes/apiserver 


#kube-controller-manager

cp /tmp/k8s_binary/kube-controller-manager /usr/bin/kube-controller-manager
cp /tmp/config/kube-controller-manager.service /usr/lib/systemd/system/kube-controller-manager.service
cp /tmp/config/controller-manager /etc/kubernetes/controller-manager 


#kube-scheduler
cp /tmp/k8s_binary/kube-scheduler /usr/bin/kube-scheduler
cp /tmp/config/kube-scheduler.service /usr/lib/systemd/system/kube-scheduler.service
cp /tmp/config/scheduler /etc/kubernetes/scheduler

#flanneld reserved config
cp /tmp/flannel/usr/bin/flanneld-start /usr/bin/flanneld-start
cp /tmp/flannel/usr/bin/flanneld /usr/bin/flanneld
mkdir /usr/libexec/flannel/
cp /tmp/flannel/usr/libexec/flannel/mk-docker-opts.sh /usr/libexec/flannel/mk-docker-opts.sh
#/usr/libexec/flannel/mk-docker-opts.sh #不需要执行
#cp /tmp/flannel/etc/sysconfig/flanneld /etc/sysconfig/flanneld  #pem versionvi /etc/sysconfig/flanneld 需要修改
cp /tmp/config/flanneld /etc/sysconfig/flanneld  
#cp /tmp/flannel/usr/lib/systemd/system/flanneld.service /usr/lib/systemd/system/flanneld.service
cp /tmp/config/flanneld.service /usr/lib/systemd/system/flanneld.service
 
    systemctl enable flanneld.service
    systemctl start flanneld.service
    systemctl status flanneld.service

#refresh service 
systemctl daemon-reload

for SERVICES in etcd kube-apiserver kube-controller-manager kube-scheduler ; do
	systemctl restart $SERVICES
	systemctl enable $SERVICES
	systemctl status $SERVICES
done


#kubectl config
cp /tmp/k8s_binary/kubectl /usr/bin/kubectl
kubectl config set-cluster default-cluster --server=http://centos-master:8080
kubectl config set-context default-context --cluster=default-cluster --user=default-admin
kubectl config use-context default-context

kubectl get nodes

kubectl cluster-info dump
