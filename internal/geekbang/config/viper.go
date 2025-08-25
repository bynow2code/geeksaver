package config

import (
	"os"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

const (
	DefaultConfigFile = "$HOME/geek-saver.yml" // 默认程序配置文件
)

var v *viper.Viper

func init() {
	v = viper.New()
	v.SetConfigFile(os.ExpandEnv(DefaultConfigFile))
}

// GetViper 返回 viper 实例
func GetViper() *viper.Viper {
	return v
}

// WriteConfig 写入程序配置文件
func WriteConfig(config *Config) error {
	var cfgMap map[string]any
	err := mapstructure.Decode(config, &cfgMap)
	if err != nil {
		return err
	}

	err = v.MergeConfigMap(cfgMap)
	if err != nil {
		return err
	}

	err = v.WriteConfig()
	if err != nil {
		return err
	}
	return nil
}
