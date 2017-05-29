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

import cmd "github.com/Huawei/containerops/singular/cmd"

func Deploynode(list map[string]string, ip string) {

	// node & ip config
	//cmd.ExecCMDparams("echo", []string{"-e", MasterIP + "centos-master\n" + NodeIP + " centos-minion", ">>", "/etc/hosts"})

	for k, v := range list {
		cmd.ExecShCommandEcho(k+" "+v, "/etc/hosts")
	}

	//docker apt-get install docker.io
	// input yes

	cmd.ExecCMDparams("apt-get", []string{"install", "docker.io"})

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
	cmd.ExecCMDparams("mkdir", []string{"-p", "/usr/lib/systemd/system/"})
	cmd.ExecCPparams("/tmp/config/kubelet.service", "/usr/lib/systemd/system/kubelet.service")
	cmd.ExecCPparams("/tmp/config/kubelet", "/etc/kubernetes/kubelet")

	//#kube-proxy
	cmd.ExecCPparams("/tmp/k8s_binary/kube-proxy", "/usr/bin/kube-proxy")
	cmd.ExecCPparams("/tmp/config/kube-proxy.service", "/usr/lib/systemd/system/kube-proxy.service")
	cmd.ExecCPparams("/tmp/config/proxy", "/etc/kubernetes/proxy")

	// #flanneld
	cmd.ExecCPparams("/tmp/flannel/usr/bin/flanneld-start", "/usr/bin/flanneld-start")

	cmd.ExecCPparams("/tmp/flannel/usr/bin/flanneld", "/usr/bin/flanneld")

	cmd.ExecCMDparams("mkdir", []string{"-p", "/usr/libexec/flannel/"})

	cmd.ExecCPparams("/tmp/flannel/usr/libexec/flannel/mk-docker-opts.sh", "/usr/libexec/flannel/mk-docker-opts.sh")
	cmd.ExecCMDparams("mkdir", []string{"-p", "/etc/sysconfig/"})
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
