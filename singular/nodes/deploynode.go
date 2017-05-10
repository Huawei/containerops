package nodes

import (
	"fmt"

	"github.com/Huawei/containerops/singular/init_config"
)

var nodelist = make(map[string]string)

func DeployNodes() {

	nodelist["centos-master"] = init_config.MasterIP
	nodelist["centos-minion"] = init_config.NodeIP
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
