package init_config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

//go get gopkg.in/yaml.v2

// cluster
var MasterIP string = "138.68.14.193"
var NodeIP string = "138.68.22.86"
var TargetIP string
var User string = "root"
var TSpet string = "6f2671de0d70ee5048379d16c0d0405df4a720ced263ffb35f67aded4834f330"
var EtcdNet string = "/kube-centos/network" // nend update config of node and master

//VM
var MSize string = "512mb"
var Region string = "sfo2"
var Slug string = "ubuntu-17-04-x64"
var Fingerprint string = "ee:81:d0:59:ab:09:1c:ff:52:dd:11:f8:bd:a6:7f:d8"

var fileslist = make(map[string]string)

func Get_files() map[string]string {

	fileslist["Etcd"] = Etcd
	fileslist["Etcdctl"] = Etcdctl
	fileslist["Flanneld"] = Flanneld
	fileslist["kube_apiserver"] = kube_apiserver_157
	fileslist["kube_controller_manager"] = kube_controller_manager_157
	fileslist["kube_proxy"] = kube_proxy_157
	fileslist["kube_scheduler"] = kube_scheduler_157
	fileslist["kubectl"] = kubectl_157
	fileslist["kubelet"] = kubelet_157

	return fileslist
}

type Settings struct {
	DBname         string `yaml:"database_name"`
	DBpass         string `yaml:"database_pass"`
	DBport         string `yaml:"database_port"`
	DBurl          string `yaml:"database_url"`
	DBuser         string `yaml:"database_user"`
	DropWhenNOrule bool   `yaml:"drop_when_no_rule"`
}

type Config struct {
	AppSettings Settings `yaml:"app_config"`
}

func getconfig(configname string) {

	filename, _ := filepath.Abs("./init.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	//
	// print everything...
	//
	fmt.Printf("%#v\n\n", config.AppSettings)
	//
	// print one by one...
	//
	fmt.Printf("database_name: %s\n", config.AppSettings.DBname)
	fmt.Printf("database_pass: %s\n", config.AppSettings.DBpass)
	fmt.Printf("database_port: %s\n", config.AppSettings.DBport)
	fmt.Printf("database_url: %s\n", config.AppSettings.DBurl)
	fmt.Printf("database_user: %s\n", config.AppSettings.DBuser)
	fmt.Printf("drop_when_no_rule: %t\n", config.AppSettings.DropWhenNOrule)
	//return config.AppSettings
}

//etcd:
var Etcd string = "https://storage.googleapis.com/containerops-release/etcd/3.1.7/etcd"
var Etcdctl string = "https://storage.googleapis.com/containerops-release/etcd/3.1.7/etcdctl"

//flannel
var Flanneld string = "https://storage.googleapis.com/containerops-release/flannel/0.7.1/flanneld"

var hyperkube string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/hyperkube"

//kubernetes 1.5.7
var kube_apiserver_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kube-apiserver"
var kube_controller_manager_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kube-controller-manager"

//var kube_discovery_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kube-discovery"
//var kube_dns_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kube-dns"
var kube_proxy_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kube-proxy"
var kube_scheduler_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kube-scheduler"
var kubectl_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kubectl"
var kubelet_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kubelet"

//var kubeadm_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kubeadm"

//var kubefed_157 string = "https://storage.googleapis.com/containerops-release/kubernetes/1.5.7/kubefed"

//kubernetes 1.6.2
var cloud_controller_manager string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/cloud-controller-manager"

//var hyperkube string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/hyperkube"
var kube_aggregator string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-aggregator"
var kube_apiserver string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-apiserver"
var kube_controller_manager string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-controller-manager"
var kube_proxy string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-proxy"
var kube_scheduler string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kube-scheduler"
var kubeadm string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kubeadm"
var kubectl string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kubectl"
var kubefed string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kubefed"
var kubelet string = "https://storage.googleapis.com/containerops-release/kubernetes/1.6.2/kubelet"
