package util

type EmailSender interface {
	SendMail(to []string, subject, body string) error
}
