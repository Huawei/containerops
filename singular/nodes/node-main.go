package nodes

import cmd "github.com/Huawei/containerops/singular/cmd"

func Deploynode(list map[string]string, ip string) {

	// node & ip config
	//cmd.ExecCMDparams("echo", []string{"-e", MasterIP + "centos-master\n" + NodeIP + " centos-minion", ">>", "/etc/hosts"})
	for k, v := range list {
		cmd.ExecShCommandEcho(k+" "+v, "/etc/hosts")
	}

	//firewalld
	cmd.ExecCMDparams("systemctl", []string{"disable", "firewalld"})
	cmd.ExecCMDparams("systemctl", []string{"stop", "firewalld"})

	//#downlload  binary & config to temp
	cmd.ExecCMDparams("mkdir", []string{"/tmp/etcd"})
	cmd.ExecCMDparams("mkdir", []string{"/tmp/k8s_binary"})

	//#erc create dir
	cmd.ExecCMDparams("mkdir", []string{"/etc/kubernetes/"})
	cmd.ExecCPparams("/tmp/config/config", "/etc/kubernetes/config")

	// #kubelet
	cmd.ExecCMDparams("mkdir", []string{"/var/lib/kubelet"})
	cmd.ExecCPparams("/tmp/k8s_binary/kubelet", "/usr/bin/kubelet")
	cmd.ExecCPparams("/tmp/config/kubelet.service", "/usr/lib/systemd/system/kubelet.service")
	cmd.ExecCPparams("/tmp/config/kubelet", "/etc/kubernetes/kubelet")

	//#kube-proxy
	cmd.ExecCPparams("/tmp/k8s_binary/kube-proxy", "/usr/bin/kube-proxy")
	cmd.ExecCPparams("/tmp/config/kube-proxy.service", "/usr/lib/systemd/system/kube-proxy.service")
	cmd.ExecCPparams("/tmp/config/proxy", "/etc/kubernetes/proxy")
	//#flanneld reserved

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

	cmd.ExecCMDparams("mkdir", []string{"/usr/libexec/flannel/"})

	cmd.ExecCPparams("/tmp/flannel/usr/libexec/flannel/mk-docker-opts.sh", "/usr/libexec/flannel/mk-docker-opts.sh")

	cmd.ExecCPparams("/tmp/config/flanneld", "/etc/sysconfig/flanneld")

	cmd.ExecCPparams("/tmp/config/flanneld.service", "/usr/lib/systemd/system/flanneld.service")

	// #refresh service
	cmd.Reload()

	// #kube-proxy kubelet  docker
	cmd.RestartSvc([]string{"kube-proxy", "kubelet", "docker"})

	// #kubectl config
	cmd.ExecCPparams("/tmp/k8s_binary/kubectl", "/usr/bin/kubectl")
	cmd.ExecCMDparams("kubectl", []string{"config", "get", "nodes"})

	// kubectl get nodes
	cmd.ExecCMDparams("kubectl", []string{"config", "cluster-info", "dump"})
}
