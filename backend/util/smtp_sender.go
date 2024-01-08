// util/smtpsender.go

package util

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type SMTPSender struct {
	Host      string
	Port      int
	Username  string
	Password  string
	UseTLS    bool
	FromEmail string
}

func (s *SMTPSender) SendMail(to []string, subject, body string) error {
	addr := fmt.Sprintf("%s:%d", s.Host, s.Port)

	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte("From: " + s.FromEmail + "\r\n" +
		"To: " + strings.Join(to, ", ") + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body)

	if s.UseTLS {
		var client *smtp.Client
		var err error

		// ポート587の場合、通常のSMTP接続を開始し、StartTLSを使用してアップグレード
		if s.Port == 587 {
			conn, err := smtp.Dial(addr)
			if err != nil {
				return err
			}
			defer conn.Close()

			tlsConfig := &tls.Config{
				ServerName:         s.Host,
				InsecureSkipVerify: false,
			}

			if err = conn.StartTLS(tlsConfig); err != nil {
				return err
			}
			client = conn
		} else if s.Port == 465 { // ポート465の場合、直接TLS接続を開始
			tlsConfig := &tls.Config{
				ServerName:         s.Host,
				InsecureSkipVerify: false,
			}

			conn, err := tls.Dial("tcp", addr, tlsConfig)
			if err != nil {
				return err
			}
			defer conn.Close()

			client, err = smtp.NewClient(conn, s.Host)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("unsupported port for TLS: %d", s.Port)
		}

		defer client.Close()

		// SMTP認証を行う
		if err = client.Auth(auth); err != nil {
			return err
		}

		// メール送信
		if err = client.Mail(s.FromEmail); err != nil {
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

	fmt.Println(addr)

	return smtp.SendMail(addr, nil, s.FromEmail, to, msg)
}
