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

package tools

const (
	//Service files folder
	ServiceFilesFolder = "service"
	KubectlFileFolder  = "kubectl"

	//Systemd server location
	SystemdServerPath = "/etc/systemd/system"
	//Default binary location
	BinaryServerPath = "/usr/local/bin"

	//Service Etcd Files Folder Name
	ServiceEtcdFolder = "etcd"
	ServiceEtcdFile   = "etcd.service"

	//Service Flannel files Folder Name
	ServiceFlanneldFolder = "flanneld"
	ServiceFlanneldFile   = "flanneld.service"

	//Service Docker files Folder Name
	ServiceDockerFolder = "docker"
	ServiceDockerFile   = "docker.service"

	//Kubectl files Folder Name
	KubectlFile                      = "kubectl"
	KubectlConfigFile                = "config"
	KubeTokenCSVFile                 = "token.csv"
	KubeBootstrapConfig              = "bootstrap.kubeconfig"
	KubeAPIServerSystemdFile         = "kube-apiserver.service"
	KubeControllerManagerSystemdFile = "kube-controller-manager.service"
	KubeSchedulerSystemdFile         = "kube-scheduler.service"
	KubeletSystemdFile               = "kubelet.service"
	KubeProxySystemdFiles            = "kube-proxy.service"

	//CA files folder
	CAFilesFolder = "ssl"

	//CA Root Files Folder Name
	CARootFilesFolder = "root"
	//CA Root Files Const Name
	CARootConfigFile    = "ca-config.json"
	CARootCSRConfigFile = "ca-csr.json"
	CARootPemFile       = "ca.pem"
	CARootCSRFile       = "ca.csr"
	CARootKeyFile       = "ca-key.pem"

	//CA Etcd Files Folder Name
	CAEtcdFolder        = "etcd"
	CAEtcdCSRConfigFile = "etcd-csr.json"
	CAEtcdKeyPemFile    = "etcd-key.pem"
	CAEtcdCSRFile       = "etcd.csr"
	CAEtcdPemFile       = "etcd.pem"

	//CA Flannel Files Folder name
	CAFlanneldFolder        = "flanneld"
	CAFlanneldCSRConfigFile = "flanneld-csr.json"
	CAFlanneldKeyPemFile    = "flanneld-key.pem"
	CAFlanneldCSRFile       = "flanneld.csr"
	CAFlanneldPemFile       = "flanneld.pem"

	//CA Docker Files Folder name
	CADockerFolder = "docker"

	//Kubernetes admin Files
	CAKubernetesFolder       = "kubernetes"
	CAKubeAdminCSRConfigFile = "admin-csr.json"
	CAKubeAdminKeyPemFile    = "admin-key.pem"
	CAKubeAdminCSRFile       = "admin.csr"
	CAKubeAdminPemFile       = "admin.pem"
	//Kubernetes API Server Files
	CAKubeAPIServerCSRConfigFile = "kubernetes-csr.json"
	CAKubeAPIServerKeyPemFile    = "kubernetes-key.pem"
	CAKubeAPIServerCSRFile       = "kubernetes.csr"
	CAKubeAPIServerPemFile       = "kubernetes.pem"
	//Kubernetes Kube Proxy Files
	CAKubeProxyServerCSRConfigFile = "kube-proxy-csr.json"
	CAKubeProxyServerKeyPemFile    = "kube-proxy-key.pem"
	CAKubeProxyServerCSR           = "kube-proxy.csr"
	CAKubeProxyServerPemFile       = "kube-proxy.pem"
	KubeProxyConfigFile            = "kube-proxy.kubeconfig"
)
