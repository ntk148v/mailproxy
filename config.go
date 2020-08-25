package main

import (
	"time"

	"github.com/spf13/viper"
)

// proxyConfig stores mail proxy related configurations.
type proxyConfig struct {
	// address - the mail proxy binding address.
	// default: `:9011`.
	// format: <host>:<port>
	address           string `yaml:"address"`
	domain            string `yaml:"domain"`
	allowInsecureAuth bool   `yaml:"allowInsecureAuth"`
	// readTimeout & writeTimeout - seconds.
	readTimeout  time.Duration `yaml:"readTimeout"`
	writeTimeout time.Duration `yaml:"writeTimeout"`
	// maxRecipients represents how many recipients mailproxy can handle.
	maxRecipients int `yaml:"maxRecipients"`
	// maxMessageBytes represents the size of message.
	maxMessageBytes int `yaml:"maxMessageBytes"`
	// retryAttempts is the number of retry attempts to send mail.
	retryAttempts int `yaml:"retryAttempts"`
	// retryDelay is the delay in seconds between consecutive retries.
	retryDelay int `yaml:"retryDelay"`
}

// smtpConfig stores smtp server related configurations.
type smtpConfig struct {
	// username - smtp username.
	username string `yaml:"username"`
	// password - smtp user's password.
	password secret `yaml:"password"`
	// address - smtp server address, format: <host>:<port>
	address string `yaml:"address"`
}

type globalConfig struct {
	proxy proxyConfig
	smtp  smtpConfig
}

// loadConfig generates a configuration instance which will be passed around the codebase
func loadConfig(cp string) error {
	viper.SetConfigFile(cp)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Set default values
	viper.SetDefault("proxy.address", ":9011")
	viper.SetDefault("proxy.domain", "localhost")
	viper.SetDefault("proxy.readTimeout", time.Duration(10))
	viper.SetDefault("proxy.writeTimeout", time.Duration(10))
	viper.SetDefault("proxy.maxRecipients", 50)
	viper.SetDefault("proxy.maxMessageBytes", 1024*1024)
	viper.SetDefault("proxy.retryAttempts", 5)
	viper.SetDefault("proxy.retryDelay", time.Duration(5))
	var cfg globalConfig
	err = viper.Unmarshal(&cfg)
	return err
}
