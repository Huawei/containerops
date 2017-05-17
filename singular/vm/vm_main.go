package vm

import "github.com/Huawei/containerops/singular/init_config"

func CreateVMs() {
	var nodelist = init_config.Get_nodes()
	for k, ip := range nodelist {
		init_config.TargetIP = ip
		//fmt.Printf("k=%v, v=%v\n", k, ip)
		if k == init_config.Master_name {
			CreateNewVM(k)
		}
		if k == init_config.Minion_name {
			CreateNode_VM(k)
		}
	}
}
