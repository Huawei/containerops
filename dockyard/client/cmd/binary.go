/*
Copyright 2016 - 2017 Huawei Technologies Co., Ltd. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Huawei/containerops/common"
	"github.com/Huawei/containerops/dockyard/client/binary"
)

var force bool

var binaryCmd = &cobra.Command{
	Use:   "binary",
	Short: "Upload or download binary file from repository build by Dockyard service",
	Long: `Use binary sub command to upload or download binary file from repository build by Dockyard service.

Upload file to a repository:

  warship binary upload --domain hub.opshub.sh /tmp/warship containerops/cncf-demo/stichers

  The upload URI pattern is <namespace>/<repository>/<tag>

Download file from repository of Dockyard:

  warship binary download --domain hub.opshub.sh containerops/cncf-demo/strichers/warship

  The download URI pattern is <namespace>/<repository>/<tag>/<filename>
`,
}

var uplaodCmd = &cobra.Command{
	Use:   "upload",
	Short: "Upload file to repository build by Dockyard service, `warship binary upload --domain hub.opshub.sh <filename> <namespace>/<repository>/<tag>`",
	Long: `Upload file to a repository of Dockyard:

warship binary upload --domain hub.opshub.sh  /tmp/warship containerops/cncf-demo/stichers

The upload URI pattern is <namespace>/<repository>/<tag>`,
	Run: uploadBinary,
}

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download file form repository build by Dockyard service, `warship binary download <namespace>/<repository>/<tag>/<filename> <path>",
	Long: `Download file from repository of Dockyard:

warship binary download --domain hub.opshub.sh  containerops/cncf-demo/stichers/warship /tmp

The download URI pattern is <namespace>/<repository>/<tag>/<filename>`,
	Run: downloadBinary,
}

//init()
func init() {
	RootCmd.AddCommand(binaryCmd)

	//Add create repository sub command.
	binaryCmd.AddCommand(uplaodCmd)
	binaryCmd.AddCommand(downloadCmd)

	uplaodCmd.Flags().BoolVarP(&force, "force", "f", false, "use the current version of file replace the file in the server")
	viper.BindPFlag("force", uplaodCmd.Flags().Lookup("force"))

}

//uploadBinary upload binary to Dockyard service.
//curl -i -X PUT -T <filename> -H "Content-Type: text/plain"  https://hub.opshub.sh/binary/v1/:namespace/:repository/binary/:tag/:binary
func uploadBinary(cmd *cobra.Command, args []string) {
	if domain == "" {
		domain = common.Warship.Domain
	}

	if len(args) <= 0 {
		fmt.Println("The file path and upload uri is required.")
		os.Exit(1)
	}

	namespace := strings.Split(args[1], "/")[0]
	repository := strings.Split(args[1], "/")[1]
	tag := strings.Split(args[1], "/")[2]

	if err := binary.UploadBinaryFile(args[0], domain, namespace, repository, tag, force); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Upload file [", args[0], "] sucessfully.")
	os.Exit(0)
}

//downloadBinary download binary file from Dockyard service.
func downloadBinary(cmd *cobra.Command, args []string) {
	if domain == "" {
		domain = common.Warship.Domain
	}

	if len(args) <= 0 {
		fmt.Println("The file path and upload uri is required.")
		os.Exit(1)
	}

	namespace := strings.Split(args[0], "/")[0]
	repository := strings.Split(args[0], "/")[1]
	tag := strings.Split(args[0], "/")[2]
	filename := strings.Split(args[0], "/")[3]

	if err := binary.DownloadBinaryFile(domain, namespace, repository, filename, tag, path.Join(args[1], filename)); err != nil {
		fmt.Println("Download file error: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Download file sucessfully.")
	os.Exit(0)
}
