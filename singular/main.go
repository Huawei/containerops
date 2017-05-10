package main

import (
	"github.com/Huawei/containerops/singular/init_config"
	"github.com/Huawei/containerops/singular/nodes"
	"github.com/Huawei/containerops/singular/vm"
)

type SSHCommander struct {
	User string
	IP   string
}

//var nodes = [2][2]string{{"192.168.60.141", "centos-master"}, {"192.168.60.150", "centos-minion"}}

func main() {
	// SSHCommander.IP
	init_config.TargetIP = init_config.MasterIP
	vm.CreateNewVM("lidian-unbantu-droplet")
	nodes.DownloadFiles()
	nodes.DeployNodes()
}
