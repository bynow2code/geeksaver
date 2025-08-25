package config

import (
	"os"
)

const DefaultMdSavePath = "$HOME/geek-docs"

// Config 程序主配置文件
type Config struct {
	User UserConfig `yaml:"user"`
	Md   MdConfig   `yaml:"md"`
}

// UserConfig 用户配置项
type UserConfig struct {
	GCID  string `mapstructure:"gcid"`
	GCESS string `mapstructure:"gcess"`
}

// MdConfig markdown配置项
type MdConfig struct {
	SavePath string `mapstructure:"savepath"`
}

// 全局配置实例
var config = &Config{}

func init() {
	config.Md.SavePath = os.ExpandEnv(DefaultMdSavePath)
}

// SetGCID 设置用户id
func SetGCID(gcid string) {
	config.User.GCID = gcid
}

// SetGCESS 设置用户令牌
func SetGCESS(gcess string) {
	config.User.GCESS = gcess
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	return config
}
