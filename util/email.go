package util

import (
	"applichic.com/chic_secret/config"
	"log"
	"net/smtp"
)

func SendEmail(email string, body string, subject string) {
	auth := smtp.PlainAuth("", config.Conf.Email, config.Conf.EmailPassword, "smtp.gmail.com")
	msg := []byte("To: " + email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail("smtp.gmail.com:587", auth, config.Conf.Email, []string{email}, msg)
	if err != nil {
		log.Fatal(err)
	}
}
