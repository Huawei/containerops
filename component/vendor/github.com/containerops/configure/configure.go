package configure

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

const (
	//ConfigName using in the config file name or environment variable prefix.
	ConfigName = "containerops"
	//ConfigFileEtcPath is the config file in the /etc path.
	ConfigFileEtcPath = "/etc/containerops"
	//configFileHomePath is the config file in the user root path, like `/home/genendna/.containerops`
	configFileHomePath = "$HOME/.containerops"
	//configFileInclude is the config file in the project path, usual in the `conf` folder.
	configFileInclude = "./conf"
)

var (
	configFile, configEnv *viper.Viper
)

func init() {
	configFile, _ = initConfigFile()
	configEnv, _ = initConfigEnv()
}

//initConfigFile is init a viper instance with file config.
func initConfigFile() (*viper.Viper, error) {
	configFile := viper.New()
	configFile.SetConfigName(ConfigName)
	configFile.SetConfigType("toml")
	configFile.AddConfigPath(ConfigFileEtcPath)
	configFile.AddConfigPath(configFileHomePath)
	configFile.AddConfigPath(configFileInclude)

	err := configFile.ReadInConfig()
	if err != nil {
		return configFile, fmt.Errorf("Fatal error config file: %s", err)
	}

	return configFile, nil
}

//initConfigEnv is init a viper instance with environment variables.
func initConfigEnv() (*viper.Viper, error) {
	configEnv := viper.New()
	configEnv.SetEnvPrefix(ConfigName)
	configEnv.AutomaticEnv()

	return configEnv, nil
}

//Get is encapsulation of viper.Get().
func Get(key string) interface{} {
	if configEnv.IsSet(key) == true {
		return configEnv.Get(key)
	}

	return configFile.Get(key)
}

//GetBool is encapsulation of viper.GetBool().
func GetBool(key string) bool {
	if configEnv.IsSet(key) == true {
		return configEnv.GetBool(key)
	}

	return configFile.GetBool(key)
}

//GetFloat64 is encapsulation of viper.GetFloat64().
func GetFloat64(key string) float64 {
	if configEnv.IsSet(key) == true {
		return configEnv.GetFloat64(key)
	}

	return configFile.GetFloat64(key)
}

//GetInt is encapsulation of viper.GetInt().
func GetInt(key string) int {
	if configEnv.IsSet(key) == true {
		return configEnv.GetInt(key)
	}

	return configFile.GetInt(key)
}

//GetString is encapsulation of viper.GetString().
func GetString(key string) string {
	if configEnv.IsSet(key) == true {
		return configEnv.GetString(key)
	}

	return configFile.GetString(key)
}

//GetStringMap is encapsulation of viper.GetStringMap().
func GetStringMap(key string) map[string]interface{} {
	if configEnv.IsSet(key) == true {
		return configEnv.GetStringMap(key)
	}

	return configFile.GetStringMap(key)
}

//GetStringMapString is encapsulation of viper.GetStringMapString().
func GetStringMapString(key string) map[string]string {
	if configEnv.IsSet(key) == true {
		return configEnv.GetStringMapString(key)
	}

	return configFile.GetStringMapString(key)
}

//GetStringSlice is encapsulation of viper.GetStringSlice().
func GetStringSlice(key string) []string {
	if configEnv.IsSet(key) == true {
		return configEnv.GetStringSlice(key)
	}

	return configFile.GetStringSlice(key)
}

//GetTime is encapsulation of viper.GetTime().
func GetTime(key string) time.Time {
	if configEnv.IsSet(key) == true {
		return configEnv.GetTime(key)
	}

	return configFile.GetTime(key)
}

//GetDuration is encapsulation of viper.GetDuration().
func GetDuration(key string) time.Duration {
	if configEnv.IsSet(key) == true {
		return configEnv.GetDuration(key)
	}

	return configFile.GetDuration(key)
}
