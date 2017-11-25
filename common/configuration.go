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

package common

import (
	"encoding/json"
	"fmt"
	"os"

	"regexp"

	homeDir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// SetConfig is setting config file path/name/type.
func SetConfig(cfgFile string) error {
	// Find home directory.
	home, err := homeDir.Dir()
	if err != nil {
		return fmt.Errorf("read $HOME envrionment error: %s", err.Error())
	}
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name "containerops" (with extension .toml).
		viper.SetConfigType("toml")
		viper.SetConfigName("containerops")
		viper.AddConfigPath("/etc/containerops/config")
		viper.AddConfigPath(fmt.Sprintf("%s/.containerops/config", home))
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("coops")
	viper.AutomaticEnv() // read in environment variables that match
	viper.WatchConfig()

	if err := viper.ReadInConfig(); err != nil {
		// If NOT found config file, automatically create it
		if ok, _ := regexp.MatchString("^Config File.*Not Found in.*", err.Error()); ok {
			fmt.Println("Not Found Config file, will auto create")

			defaultPath := fmt.Sprintf("%s/.containerops/config", home)
			if _, err := os.Stat(defaultPath); err != nil {
				if os.IsNotExist(err) {
					if err := os.MkdirAll(defaultPath, 0777); err != nil {
						return nil
					}
				} else {
					return err
				}
			}

			newConfigFile, err := os.Create(defaultPath + "/containerops.toml")
			if err != nil {
				return fmt.Errorf("Automatically create config file ERROR: %s", err.Error())
			}
			defer newConfigFile.Close()
		} else {
			return fmt.Errorf("fatal error config file: %s", err.Error())
		}
	}

	if err := setDatabaseConfig(viper.GetStringMap("database")); err != nil {
		return err
	}

	if err := setWebConfig(viper.GetStringMap("web")); err != nil {
		return err
	}

	if err := setStorageConfig(viper.GetStringMap("storage")); err != nil {
		return err
	}

	if err := setWarshipConfig(viper.GetStringMap("warship")); err != nil {
		return err
	}

	if err := setSingularConfig(viper.GetStringMap("singular")); err != nil {
		return err
	}

	if err := setAssemblingConfig(viper.GetStringMap("assembling")); err != nil {
		return err
	}

	if err := setMailConfig(viper.GetStringMap("mail")); err != nil {
		return err
	}

	return nil
}

/*
Configurations for all modules

# 1. Configurations of database.

[database]
driver = "mysql"
host = "127.0.0.1"
port = 3306
user = "root"
password = "containerops_database"
db = "containerops_password"

# 2. Configurations for HTTPS or Unix Socket
   2.1 If multi modules deploy in one node, there should have a proxy like Caddy or Nginx.
       Each module use with Unix Socket type,  configurations look like this:

           [web]
           mode = "unix"
           address = "/var/run/${module}.socket"

   2.2 If module deploys in one node alone, it only supports HTTPS model and must have the SSL
       certification files.

[web]
domain = "opshub.sh"
mode = "https"
address = "127.0.0.1"
port = 443
cert = "PATH_TO_CERT_FILE"
key = "PATH_TO_KEY_FILE"

# 3. Configurations for storage path of Dockyard module.
#   3.1 TODO Using the Object Storage Service in the Dockyard module.

[storage]
dockerv2 = "/tmp/dockerv2" # path for image files of Docker Distribution V2 Protocol
binaryv1 = "/tmp/binaryv1" # path for binary files of Dockyard Binary V1 Protocol

# 4. Configurations for Warship of Dockyard client.

[warship]
domain = "hub.opshub.sh"

# 5. Configurations for Singular modules.

[singular]
provider = "digitalocean"
token = "435a054fba66cb11d6b7abeaa3d89aac777d4d1d"
*/

type DatabaseConfig struct {
	Driver   string `json:"driver" yaml:"driver"`
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	Name     string `json:"db" yaml:"db"`
}

type WebConfig struct {
	Domain  string `json:"domain" yaml:"domain"`
	Mode    string `json:"mode" yaml:"mode"`
	Address string `json:"address" yaml:"address"`
	Port    int    `json:"port" yaml:"port"`
	Key     string `json:"key" yaml:"key"`
	Cert    string `json:"cert" yaml:"cert"`
}

type StorageConfig struct {
	DockerV2 string `json:"dockerv2" yaml:"dockerv2"`
	BinaryV1 string `json:"binaryv1" yaml:"binaryv1"`
}

type WarshipConfig struct {
	Domain string
}

type SingularConfig struct {
	Provider string `json:"provider" yaml:"provider"`
	Token    string `json:"token" yaml:"token"`
}

type AssemblingConfig struct {
	Domain            string `json:"domain" description:"Listen domain, the official domain is *.osphub.sh"`
	Mode              string `json:"mode" description:"Listen mode, 'https' or 'unix'"`
	Address           string `json:"address" description:"The host address when mode is 'https', or socket file path when mode is 'unix'"`
	Port              int    `json:"port"`
	Key               string `json:"key"`
	Cert              string `json:"cert"`
	DockerDaemonImage string `json:"docker_daemon_image" description:"Image with a docker daemon, providing Docker Engine APIs"`
	KubeConfig        string `json:"kubeconfig" description:"The address of k8s api server"`
	ServiceType       string `json:"service_type" description:"The service type of the dind environment, might be 'NodePort' or 'LoadBalancer'"`
}

type MailConfig struct {
	SmtpAddress string `json:"smtp_address" yaml:"smtp_address"`
	SmtpPort    string `json:"smtp_port" yaml:"smtp_port"`
	User        string `json:"user" yaml:"user"`
	Password    string `json:"password" yaml:"password"`
}

var Database DatabaseConfig
var Web WebConfig
var Storage StorageConfig
var Warship WarshipConfig
var Singular SingularConfig
var Assembling AssemblingConfig
var Mail MailConfig

func setDatabaseConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &Database)
}

func setWebConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &Web)
}

func setStorageConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &Storage)
}

func setWarshipConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &Warship)
}

func setSingularConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	return json.Unmarshal(bs, &Singular)
}

func setAssemblingConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &Assembling)
	if err != nil {
		return err
	}

	return nil
}

func setMailConfig(config map[string]interface{}) error {
	bs, err := json.Marshal(&config)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bs, &Mail)
	if err != nil {
		return err
	}
	return nil
}
