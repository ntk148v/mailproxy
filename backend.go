package main

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/ntk148v/go-smtp"
	"github.com/spf13/viper"
)

// The Backend implements SMTP server methods.
type Backend struct{}

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
	if (username != smtpusername && username != strings.Split(smtpusername, "@")[0] && strings.Split(username, "@")[0] != smtpusername) || password != viper.GetString("smtp.password") {
		return nil, errors.New("Invalid username or password")
	}
	log.Println("---------------")
	return &Session{}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (be *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

// Mail - handle MAIL FROM
func (s *Session) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	s.Reset()
	s.msg.From = from
	return nil
}

// Rcpt - handle RCPT TO
func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	s.msg.To = append(s.msg.To, to)
	return nil
}

// Data - handle DATA
func (s *Session) Data(r io.Reader) error {
	auth := sasl.NewPlainClient("", viper.GetString("smtp.username"), viper.GetString("smtp.password"))
	err := smtp.SendMail(viper.GetString("smtp.address"), auth, s.msg.From, s.msg.To, r, true)
	if err != nil {
		log.Println("Error when handle data: ", err)
		return err
	}
	log.Println("Nice! Forwarded mail to ", s.msg.To)
	return nil
}

// Reset - Clear message
func (s *Session) Reset() {
	s.msg = &message{}
}

// Logout - logout, of course
func (s *Session) Logout() error {
	log.Println("---------------")
	return nil
}
