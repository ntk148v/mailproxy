package main

import (
	"errors"
	"io"
	"strings"
	"time"

	"github.com/emersion/go-sasl"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ntk148v/go-smtp"
	"github.com/spf13/viper"
)

// The Backend implements SMTP server methods.
type Backend struct {
	logger log.Logger
}

type message struct {
	From string
	To   []string
}

// A Session is returned after successful login.
type Session struct {
	backend *Backend
	msg     *message
}

// Login handles a login command with username and password.
func (be *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	smtpusername := viper.GetString("smtp.username")
	if (username != smtpusername && username != strings.Split(smtpusername, "@")[0] &&
		strings.Split(username, "@")[0] != smtpusername) || password != viper.GetString("smtp.password") {
		return nil, errors.New("Invalid username or password")
	}
	level.Info(be.logger).Log("msg", "===================================================")
	level.Info(be.logger).Log("msg", "handle login with username and password")
	return &Session{backend: be}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (be *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

// Mail - handle MAIL FROM
func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	s.Reset()
	level.Info(s.backend.logger).Log("msg", "handle smtp command MAIL FROM", "sender", from)
	s.msg.From = from
	return nil
}

// Rcpt - handle RCPT TO
func (s *Session) Rcpt(to string) error {
	level.Info(s.backend.logger).Log("msg", "handle smtp command RCPT TO", "recipient", to)
	s.msg.To = append(s.msg.To, to)
	return nil
}

// Data - handle DATA
func (s *Session) Data(r io.Reader) error {
	level.Info(s.backend.logger).Log("msg", "handle smtp command DATA")
	auth := sasl.NewPlainClient("", viper.GetString("smtp.username"), viper.GetString("smtp.password"))
	retry := 0
	for retry < viper.GetInt("proxy.retryAttempts") {
		err := smtp.SendMail(viper.GetString("smtp.address"), auth, s.msg.From, s.msg.To, r, true)
		if err != nil {
			level.Error(s.backend.logger).Log("msg", "error when handling data", "err", err.Error())
			// Only retry if there is 4xx smtp error.
			// For details, please check: https://serversmtp.com/smtp-error/
			if smtpErr, ok := err.(*smtp.SMTPError); ok && smtpErr.Temporary() {
				retry++
				level.Debug(s.backend.logger).Log("msg", "retry to send mail", "attempt", retry)
				time.Sleep(viper.GetDuration("proxy.retryDelay") * time.Second)
				continue
			}
			return err
		}
		break
	}
	level.Info(s.backend.logger).Log("msg", "forwared mail to", "recipients", strings.Join(s.msg.To, ","))
	return nil
}

// Reset - Clear message
func (s *Session) Reset() {
	s.msg = &message{}
}

// Logout - logout, of course
func (s *Session) Logout() error {
	level.Info(s.backend.logger).Log("msg", "handle Logout")
	level.Info(s.backend.logger).Log("msg", "===================================================")
	return nil
}
