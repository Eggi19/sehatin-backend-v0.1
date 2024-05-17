package utils

import (
	"bytes"
	"fmt"
	"net/smtp"
	"text/template"
)

type EmailSender interface {
	SendEmail(email string, message string, subject string) error
}

type GoogleEmailSender struct {
}

func NewGoogleEmailSender() *GoogleEmailSender {
	return &GoogleEmailSender{}
}

func (e *GoogleEmailSender) SendEmail(email string, message string, subject string) error {
	config, err := ConfigInit()
	if err != nil {
		return err
	}

	from := config.Email
	password := config.EmailPassword

	to := []string{
		email,
	}

	smtpHost := config.SmtpHost
	smtpPort := config.SmtpPort

	auth := smtp.PlainAuth("", from, password, smtpHost)

	t, _ := template.ParseFiles("./utils/template.html")

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: %s \n%s\n\n", subject, mimeHeaders)))

	t.Execute(&body, struct {
		Message string
	}{
		Message: message,
	})

	err = smtp.SendMail(fmt.Sprintf("%s:%s", smtpHost, smtpPort), auth, from, to, body.Bytes())
	if err != nil {
		return err
	}

	return nil
}
