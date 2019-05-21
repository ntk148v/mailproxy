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
	// The path to Cert and Key files.
	// For example, to generate those files:
	// # Key considerations for algorithm "RSA" ≥ 2048-bit
	// openssl genrsa -out server.key 2048
	//
	// # Key considerations for algorithm "ECDSA" (X25519 || ≥ secp384r1)
	// # https://safecurves.cr.yp.to/
	// # List ECDSA the supported curves (openssl ecparam -list_curves)
	// openssl ecparam -genkey -name secp384r1 -out server.key
	// openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	serverCrt string `yaml:"serverCrt"`
	serverKey string `yaml:"serverKey"`
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
	viper.SetConfigName("config")
	viper.AddConfigPath(cp)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Set default values
	viper.SetDefault("proxy.address", ":9011")
	viper.SetDefault("proxy.domain", "localhost")
	viper.SetDefault("proxy.readTimeout", time.Duration(10))
	viper.SetDefault("proxy.writeTimeout", time.Duration(10))
	// By default, put cert files in the config directory
	viper.SetDefault("proxy.serverCrt", cp+"/server.crt")
	viper.SetDefault("proxy.serverKey", cp+"/server.key")
	viper.SetDefault("maxRecipients", 50)
	viper.SetDefault("maxMessageBytes", 1024*1024)
	var cfg globalConfig
	err = viper.Unmarshal(&cfg)
	return err
}
