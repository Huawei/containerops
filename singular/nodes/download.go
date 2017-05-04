package nodes

import (
	"fmt"

	cmd "github.com/Huawei/containerops/singular/cmd"

	init_config "github.com/Huawei/containerops/singular/init_config"
)

func DownloadFiles() {

	var fileslist = init_config.Get_files()
	for key, url := range fileslist {
		cmd.ExecCMDparams("wget", []string{"-c", url})
		fmt.Printf("%s\n\n", key)

	}
}
