package nodes

import (
	"fmt"

	cmd "github.com/Huawei/containerops/singular/cmd"

	init_config "github.com/Huawei/containerops/singular/init_config"
)

func DownloadFiles() {

	var fileslist = init_config.Get_files()
	for key, url := range fileslist {
		cmd.ExecCMDparams("wget", []string{"-P", "/tmp/", "-c", url})
		fmt.Printf("%s\n\n", key)

	}
	//scp -r ./config/. root@138.68.14.193:/tmp/
	cmd.LocalExecCMDparams("scp", []string{"-r", "./config/.", init_config.User + "@" + init_config.TargetIP + ":" + "/tmp/config/"})

}
