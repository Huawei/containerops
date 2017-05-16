package main

import (
	"github.com/Huawei/containerops/singular/deploy"
	"github.com/Huawei/containerops/singular/download"
	"github.com/Huawei/containerops/singular/vm"
)

type SSHCommander struct {
	User string
	IP   string
}

//var nodes = [2][2]string{{"192.168.60.141", "centos-master"}, {"192.168.60.150", "centos-minion"}}

func main() {
	// SSHCommander.IP
	//init_config.TargetIP = init_config.MasterIP

	// create vmlist
	vm.CreateNewVM("lidian-unbantu-wk-master")
	//get iplist while
	download.Download_main()
	deploy.DeployNodes()

}
