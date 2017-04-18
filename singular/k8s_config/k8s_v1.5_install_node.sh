

#config your IP
echo "192.168.121.9	centos-master
192.168.121.65	centos-minion" >> /etc/hosts

#firewalld
systemctl disable firewalld
systemctl stop firewalld

#change to binary way
yum install docker

#downlload  binary & config to temp
/tmp/k8s_binary

#create dir 
	mkdir /etc/kubernetes/
	cp /tmp/config/config /etc/kubernetes/config

	
#kubelet
	mkdir /var/lib/kubelet
	cp /tmp/k8s_binary/kubelet /usr/bin/kubelet
	cp /tmp/config/kubelet.servic /usr/lib/systemd/system/kubelet.service
	cp /tmp/config/kubelet /etc/kubernetes/kubelet 
	#/etc/kubernetes/config


#Kube-proxy

cp /tmp/k8s_binary/Kube-proxy /usr/bin/Kube-proxy
cp /tmp/config/kube-proxy.service /usr/lib/systemd/system/kube-proxy.service
cp /tmp/config/proxy /etc/kubernetes/proxy 


#flanneld reserved 
#  cp /tmp/flannel/usr/bin/flanneld-start /usr/bin/flanneld-start
cp /tmp/flannel/usr/bin/flanneld /usr/bin/flanneld
mkdir /usr/libexec/flannel/
cp /tmp/flannel/usr/libexec/flannel/mk-docker-opts.sh /usr/libexec/flannel/mk-docker-opts.sh
/usr/libexec/flannel/mk-docker-opts.sh
cp /tmp/flannel/etc/sysconfig/flanneld /etc/sysconfig/flanneld 
cp /tmp/flannel/usr/lib/systemd/system/flanneld.service /usr/lib/systemd/system/flanneld.service
 
 FLANNEL_ETCD_ENDPOINTS="http://centos-master:2379"
 FLANNEL_ETCD_PREFIX="/kube-centos/network"

    systemctl enable flanneld.service
    systemctl start flanneld.service
    systemctl status flanneld.service
#reload

systemctl daemon-reload

for SERVICES in kube-proxy kubelet  docker; do
    systemctl restart $SERVICES
    systemctl enable $SERVICES
    systemctl status $SERVICES
done

