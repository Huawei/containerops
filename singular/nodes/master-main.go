package nodes

import cmd "github.com/Huawei/containerops/singular/cmd"
import init_config "github.com/Huawei/containerops/singular/init_config"

func Deploymaster(list map[string]string, ip string) {

	// node & ip config
	//echo -e "192.168.121.9 centos-master\n192.168.121.65 centos-minion" >> /etc/hosts
	//	cmd.ExecCMDparams("echo", []string{"-e", MasterIP + "centos-master\n" + NodeIP + " centos-minion", ">>", "/etc/hosts"})

	for k, v := range list {
		cmd.ExecShCommandEcho(k+" "+v, "/etc/hosts")
	}

	// firewalld
	cmd.ExecCMDparams("systemctl", []string{"disable", "firewalld"})
	cmd.ExecCMDparams("systemctl", []string{"stop", "firewalld"})
	//#downlload  binary & config to temp
	cmd.ExecCMDparams("mkdir", []string{"/tmp/etcd"})
	cmd.ExecCMDparams("mkdir", []string{"/tmp/k8s_binary"})

	//#etc create dir
	cmd.ExecCMDparams("mkdir", []string{"/etc/kubernetes/"})
	cmd.ExecCPparams("/tmp/config/config", "/etc/kubernetes/config")
	// #etc service
	cmd.ExecCMDparams("mkdir", []string{"/etc/etcd/"})
	cmd.ExecCPparams("/tmp/etcd/etc/etcd/etcd.conf", "/etc/etcd/etcd.conf")
	cmd.ExecCMDparams("mkdir", []string{"/usr/lib/systemd/system/"})
	cmd.ExecCPparams("/tmp/config/etcd.service", "/usr/lib/systemd/system/etcd.service")
	cmd.ExecCPparams("/tmp/etcd/usr/bin/etcd", "/usr/bin/etcd")
	cmd.ExecCPparams("/tmp/etcd/usr/bin/etcdctl", "/usr/bin/etcdctl")
	cmd.ExecCMDparams("mkdir", []string{"/var/lib/etcd/"})
	// #/var/lib/etcd/default.etcd auto?

	cmd.ExecCMDparams("systemctl", []string{"daemon-reload"})
	cmd.ServiceStart("etcd")
	cmd.ServiceIsEnabled("etcd")
	cmd.ServiceExists("etcd")
	cmd.ExecCMDparams("etcdctl", []string{"mkdir", init_config.EtcdNet}) //"/kube-centos/network/config nend update config
	cmd.ExecCMDparams("etcdctl", []string{"mk", init_config.EtcdNet, "{\"Network\":\"172.40.0.0/16\"\\,\"SubnetLen\":24\\,\"Backend\":{\"Type\":\"vxlan\"}}"})

	//#kube-apiserver
	cmd.ExecCPparams("/tmp/k8s_binary/kube-apiserver", "/usr/bin/kube-apiserver")
	cmd.ExecCPparams("/tmp/config/kube-apiserver.service", "/usr/lib/systemd/system/kube-apiserver.service")
	cmd.ExecCPparams("/tmp/config/apiserver", "/etc/kubernetes/apiserver")

	//#kube-controller-manager
	cmd.ExecCPparams("/tmp/k8s_binary/kube-controller-manager", "/usr/bin/kube-controller-manager")
	cmd.ExecCPparams("/tmp/config/kube-controller-manager.service", "/usr/lib/systemd/system/kube-controller-manager.service")
	cmd.ExecCPparams("/tmp/config/controller-manager", "/etc/kubernetes/controller-manager")
	//#kube-scheduler
	cmd.ExecCPparams("/tmp/k8s_binary/kube-scheduler", "/usr/bin/kube-scheduler")
	cmd.ExecCPparams("/tmp/config/kube-scheduler.service", "/usr/lib/systemd/system/kube-scheduler.service")
	cmd.ExecCPparams("/tmp/config/scheduler", "/etc/kubernetes/scheduler")

	// #flanneld reserved
	cmd.ExecCPparams("/tmp/flannel/usr/bin/flanneld-start", "/usr/bin/flanneld-start")

	cmd.ExecCPparams("/tmp/flannel/usr/bin/flanneld", "/usr/bin/flanneld")

	cmd.ExecCMDparams("mkdir", []string{"-p", "/usr/libexec/flannel/"})
	cmd.ExecCMDparams("mkdir", []string{"/tmp/etcd"})
	// check work
	cmd.ExecCPparams("/tmp/flannel/usr/libexec/flannel/mk-docker-opts.sh", "/usr/libexec/flannel/mk-docker-opts.sh")
	cmd.ExecCMDparams("mkdir", []string{"-p", "/etc/sysconfig/"})

	cmd.ExecCPparams("/tmp/config/flanneld", "/etc/sysconfig/flanneld")

	cmd.ExecCPparams("/tmp/config/flanneld.service", "/usr/lib/systemd/system/flanneld.service")

	// #refresh service
	cmd.Reload()

	// #for SERVICES in etcd kube-apiserver kube-controller-manager kube-scheduler ; do
	// 	#systemctl restart $SERVICES
	// 	#systemctl enable $SERVICES
	// 	#systemctl status $SERVICES
	// #done
	cmd.RestartSvc([]string{"etcd", "kube-apiserver", "kube-controller-manager", "kube-scheduler"})

	// #kubectl config
	cmd.ExecCPparams("/tmp/k8s_binary/kubectl", "/usr/bin/kubectl")
	cmd.ExecCMDparams("kubectl", []string{"config", "set-cluster", "default-cluster", "--server=http://centos-master:8080"})
	cmd.ExecCMDparams("kubectl", []string{"config", "set-cluster", "default-cluster", "--cluster=default-cluster", "--user=default-admin"})
	cmd.ExecCMDparams("kubectl", []string{"config", "use-context", "use-context"})
	cmd.ExecCMDparams("kubectl", []string{"config", "get", "nodes"})

	// // kubectl get nodes
	cmd.ExecCMDparams("kubectl", []string{"config", "cluster-info", "dump"})

}
