package main

import (
	"errors"
	"io"
	"log"

	"github.com/emersion/go-sasl"
	"github.com/ntk148v/go-smtp"
	"github.com/spf13/viper"
)

// The Backend implements SMTP server methods.
type Backend struct{}

type message struct {
	From string
	To   []string
	Data []byte
}

// A Session is returned after successful login.
type Session struct {
	backend *Backend
	msg     *message
}

// Login handles a login command with username and password.
func (be *Backend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	if username != viper.GetString("smtp.username") || password != viper.GetString("smtp.password") {
		return nil, errors.New("Invalid username or password")
	}
	return &Session{}, nil
}

// AnonymousLogin requires clients to authenticate using SMTP AUTH before sending emails
func (be *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return nil, smtp.ErrAuthRequired
}

func (s *Session) Mail(from string) error {
	log.Println("Mail from:", from)
	s.Reset()
	s.msg.From = from
	return nil
}

func (s *Session) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	s.msg.To = append(s.msg.To, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	auth := sasl.NewPlainClient("", viper.GetString("smtp.username"), viper.GetString("smtp.password"))
	err := smtp.SendMail(viper.GetString("smtp.address"), auth, s.msg.From, s.msg.To, r, true)
	if err != nil {
		return err
	}
	return nil
}

func (s *Session) Reset() {
	s.msg = &message{}
}

func (s *Session) Logout() error {
	return nil
}
