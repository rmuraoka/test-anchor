// util/smtpsender.go

package util

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPSender struct {
	Host     string
	Port     int
	Username string
	Password string
	UseTLS   bool
}

func (s *SMTPSender) SendMail(to []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte("From: " + s.Username + "\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body)

	// TLSを使用する場合
	if s.UseTLS {
		tlsConfig := &tls.Config{
			ServerName:         s.Host,
			InsecureSkipVerify: false,
		}

		conn, err := tls.Dial("tcp", addr, tlsConfig)
		if err != nil {
			return err
		}
		client, err := smtp.NewClient(conn, s.Host)
		if err != nil {
			return err
		}
		defer client.Close()

		// SMTP認証を行う
		if err = client.Auth(auth); err != nil {
			return err
		}

		// メール送信
		if err = client.Mail(s.Username); err != nil {
			return err
		}
		if err = client.Rcpt(strings.Join(to, ",")); err != nil {
			return err
		}
		w, err := client.Data()
		if err != nil {
			return err
		}
		_, err = w.Write(msg)
		if err != nil {
			return err
		}
		err = w.Close()
		if err != nil {
			return err
		}
		return client.Quit()
	}

	// TLSを使用しない場合の通常のメール送信
	return s.sendMailWithoutTLS(to, subject, body, addr)
}

func (s *SMTPSender) sendMailWithoutTLS(to []string, subject, body string, addr string) error {
	msg := []byte(fmt.Sprintf("To: %s\nSubject: %s\n\n%s", strings.Join(to, ","), subject, body))

	return smtp.SendMail(addr, nil, s.Username, to, msg)
}
