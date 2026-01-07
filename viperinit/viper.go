package viperinit

import (
	"strings"

	"github.com/spf13/viper"
)

type Option func(*viper.Viper)

func NewViper(name, configType, path string, opts ...Option) *viper.Viper {
	v := viper.New()

	if name == "" {
		name = "config"
	}
	if configType == "" {
		configType = "yaml"
	}
	if path == "" {
		path = "."
	}

	// set config file name and type
	v.SetConfigName(name)
	v.SetConfigType(configType)
	v.AddConfigPath(path)

	// bind environment variables
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// apply options
	for _, opt := range opts {
		opt(v)
	}

	// read config file
	err := v.ReadInConfig()
	if err != nil {
		panic("failed to read config file: " + err.Error())
	}
	return v
}
