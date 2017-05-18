package download

import "github.com/Huawei/containerops/singular/init_config"

func Download_main() {
	var nodelist = init_config.Get_nodes()
	for k, ip := range nodelist {
		init_config.TargetIP = ip
		//fmt.Printf("k=%v, v=%v\n", k, ip)
		if k == init_config.Master_name {
			Master_Download(ip)
		}
		if k == init_config.Minion_name {
			Node_download(ip)
		}
	}
}
