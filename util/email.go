package util

import (
	"applichic.com/chic_secret/config"
	"log"
	"net/smtp"
)

// SendEmail Send an email from Chic Secret to the user
func SendEmail(email string, body string, subject string) {
	auth := smtp.PlainAuth("Chic Secret", config.Conf.Email, config.Conf.EmailPassword, "ssl0.ovh.net")
	msg := []byte("To: " + email + "\r\n" +
		"From: " + config.Conf.Email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail("ssl0.ovh.net:587", auth, config.Conf.Email, []string{email}, msg)
	if err != nil {
		log.Fatal(err)
	}
}
