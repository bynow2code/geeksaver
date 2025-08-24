package geekbang

import (
	"log"
	"os"
	"sync"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

const (
	CONFIG_PATH = "$HOME/geekbang-doc-saver.yml"
)

type ViperSingleton struct {
	once  sync.Once
	viper *viper.Viper
}

var defaultViperSingleton = &ViperSingleton{}

func GetViper() *ViperSingleton {
	defaultViperSingleton.once.Do(func() {
		viper.SetConfigFile(os.ExpandEnv(CONFIG_PATH))
		defaultViperSingleton.viper = viper.GetViper()
	})
	return defaultViperSingleton
}

type Config struct {
	User UserConfig `yaml:"user"`
}

type UserConfig struct {
	GCID  string `mapstructure:"gcid"`
	GCESS string `mapstructure:"gcess"`
}

func (v *ViperSingleton) WriteConfig(config *Config) {
	var cfgMap map[string]any
	err := mapstructure.Decode(config, &cfgMap)
	if err != nil {
		log.Fatalf("Error decoding config: %v", err)
	}
	err = v.viper.MergeConfigMap(cfgMap)
	if err != nil {
		log.Fatalf("Error merging config: %v", err)
	}
	err = v.viper.WriteConfig()
	if err != nil {
		log.Fatalf("write config failed: %v", err)
	}
}
