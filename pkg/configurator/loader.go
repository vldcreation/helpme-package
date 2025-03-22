package configurator

import (
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func LoadFromYaml(fpath string, cfg interface{}) error {
	b, err := os.ReadFile(fpath)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, cfg)
}

func LoadFromYamlOrPanic(fpath string, cfg interface{}) {
	if err := LoadFromYaml(fpath, cfg); err != nil {
		panic(err)
	}
}

func LoadFromEnv(fpath string, cfg interface{}) error {
	viper.SetConfigFile(fpath)
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	return viper.Unmarshal(cfg)
}
