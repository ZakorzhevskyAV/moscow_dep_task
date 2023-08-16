package config

import (
	"gopkg.in/yaml.v3"
	"moscow_dep_task/types"
	"os"
)

func ParseConfig(path string, config *types.Config) error {
	configbytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(configbytes, config)
	if err != nil {
		return err
	}
	return err
}
