package main

import "github.com/spf13/viper"

type proxyConfig struct {
	address string `yaml:"address"`
	domain  string `yaml:"domain"`
}

type smtpConfig struct {
	username string `yaml:"username"`
	password secret `yaml:"password"`
	address  string `yaml:"address"`
}

type globalConfig struct {
	proxy proxyConfig
	smtp  smtpConfig
}

// loadConfig generates a configuration instance which will be passed around the codebase
func loadConfig(cp string) error {
	viper.SetConfigName("config")
	viper.AddConfigPath(cp)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Set default values
	viper.SetDefault("proxy.address", ":9011")
	viper.SetDefault("proxy.domain", "localhost")
	var cfg globalConfig
	err = viper.Unmarshal(&cfg)
	return err
}
