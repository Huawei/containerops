package main

import (
	"github.com/Huawei/containerops/singular/deploy"
	"github.com/Huawei/containerops/singular/download"
	"github.com/Huawei/containerops/singular/vm"
)

func main() {
	// SSHCommander.IP
	//init_config.TargetIP = init_config.MasterIP

	// create vmlist
	vm.CreateVMs()
	//vm.CreateNewVM("lidian-unbantu-wk-master")
	//get iplist while
	download.Download_main()
	deploy.DeployNodes()

}
