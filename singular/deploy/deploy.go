package deploy

import (
	"fmt"

	"github.com/Huawei/containerops/singular/init_config"
)

func DeployNodes() {
	var nodelist = init_config.Get_nodes()
	for k, ip := range nodelist {
		init_config.TargetIP = ip
		fmt.Printf("k=%v, v=%v\n", k, ip)
		if k == "centos-master" {
			Deploymaster(nodelist, ip)
		}
		if k == "centos-minion" {
			Deploynode(nodelist, ip)
		}
	}
}
