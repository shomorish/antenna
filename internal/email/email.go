package email

import (
	"net/smtp"
)

func SendEmail(from string, password string, smtpHost string, to string, message []byte) error {
	auth := smtp.PlainAuth(
		"",
		from,
		password,
		smtpHost,
	)

	return smtp.SendMail(
		smtpHost+":587",
		auth,
		from,
		[]string{to},
		message,
	)
}

func SendEmails(from string, password string, smtpHost string, to []string, message []byte) error {
	for _, t := range to {
		if err := SendEmail(from, password, smtpHost, t, message); err != nil {
			return err
		}
	}
	return nil
}
