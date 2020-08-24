package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ntk148v/go-smtp"
	"github.com/pkg/errors"
	"github.com/prometheus/common/promlog"
	logflag "github.com/prometheus/common/promlog/flag"
	"github.com/spf13/viper"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	if os.Getenv("DEBUG") != "" {
		runtime.SetBlockProfileRate(20)
		runtime.SetMutexProfileFraction(20)
	}

	cfg := struct {
		configFile string
		logConfig  promlog.Config
	}{
		logConfig: promlog.Config{},
	}

	a := kingpin.New(filepath.Base(os.Args[0]), "The mailproxy")
	a.HelpFlag.Short('h')
	a.Flag("config.file", "Mailproxy configuration file path.").
		Default("/etc/mailproxy/config.yml").StringVar(&cfg.configFile)
	logflag.AddFlags(a, &cfg.logConfig)
	_, err := a.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing commandline arguments"))
		a.Usage(os.Args[1:])
		os.Exit(2)
	}

	logger := promlog.New(&cfg.logConfig)
	level.Info(logger).Log("msg", "Staring mailproxy")

	// Load configs
	if err := loadConfig(cfg.configFile); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing config file"))
		os.Exit(2)
	}

	be := &Backend{
		logger: log.With(logger, "component", "SMTP backend"),
	}
	s := smtp.NewServer(be)

	// Generate a fake cert
	cer, err := tls.LoadX509KeyPair(viper.GetString("proxy.serverCrt"),
		viper.GetString("proxy.serverKey"))
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrapf(err, "Error parsing cert files"))
		os.Exit(2)
	}
	s.TLSConfig = &tls.Config{
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cer},
	}
	s.Addr = viper.GetString("proxy.address")
	s.Domain = viper.GetString("proxy.domain")
	s.ReadTimeout = viper.GetDuration("proxy.readTimeout") * time.Second
	s.WriteTimeout = viper.GetDuration("proxy.writeTimeout") * time.Second
	s.MaxMessageBytes = viper.GetInt("proxy.maxMessageBytes")
	s.MaxRecipients = viper.GetInt("proxy.maxRecipients")
	s.AllowInsecureAuth = viper.GetBool("proxy.allowInsecureAuth")
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

	var (
		term = make(chan os.Signal, 1)
		srvc = make(chan struct{})
	)

	go func() {
		level.Info(logger).Log("msg", "Listening", "address", s.Addr)
		if err := s.ListenAndServe(); err != nil {
			level.Error(logger).Log("msg", "Listen error", "err", err)
			close(srvc)
		}
	}()
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-term:
			level.Info(logger).Log("msg", "Received SIGTERM, exiting gracefully...")
			os.Exit(0)
		case <-srvc:
			os.Exit(1)
		}
	}
}
