package main

import (
	"flag"
	"log"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/ntk148v/go-smtp"
	"github.com/spf13/viper"
)

const (
	defaultConfigFilePath = "./etc/"
	configFilePathUsage   = "config file directory (eg. '/etc/mailproxy/'). Config file must be named 'config.yml'."
)

var configFilePath string

func init() {
	flag.StringVar(&configFilePath, "conf", defaultConfigFilePath, configFilePathUsage)
	flag.Parse()
	loadConfig(configFilePath)
}

func main() {
	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = viper.GetString("proxy.address")
	s.Domain = viper.GetString("proxy.domain")
	// TODO(kiennt): Add config later
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true
	// Force use AUTH LOGIN instead of AUTH PLAIN
	s.EnableAuth(sasl.Login, func(conn *smtp.Conn) sasl.Server {
		return sasl.NewLoginServer(func(username, password string) error {
			state := conn.State()
			session, err := be.Login(&state, username, password)
			if err != nil {
				return err
			}

			conn.SetSession(session)
			return nil
		})
	})

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
