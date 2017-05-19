package main

import "github.com/Huawei/containerops/singular/cmd"

/

func main() {

	cmd.Execute()
	// create vmlist
	vm.CreateVMs()
	//vm.CreateNewVM("lidian-unbantu-wk-master")
	//get iplist while
	download.Download_main()
	deploy.DeployNodes()

}
