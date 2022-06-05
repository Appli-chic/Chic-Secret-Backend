package util

import (
	"applichic.com/chic_secret/config"
	"log"
	"net/smtp"
)

// SendEmail Send an email from Chic Secret to the user
func SendEmail(email string, code string, subject string) {
	auth := smtp.PlainAuth("Chic Secret", config.Conf.Email, config.Conf.EmailPassword, "ssl0.ovh.net")

	body := "" +
		"<div style=\"background-color: #0B84FF; color: white;font-size: 32px;font-weight: bold;padding-left: 50px;padding-top: 16px; padding-bottom: 16px\">Chic Secret</div>" +
		"<div style=\"background-color: #001529;color: white;font-size: 16px;padding-left: 50px;padding-top: 16px; padding-bottom: 16px;\">Code: " + code + "</div>" +
		""

	msg := []byte("To: " + email + "\r\n" +
		"From: " + config.Conf.Email + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		"\r\n" + body + "\r\n")

	err := smtp.SendMail("ssl0.ovh.net:587", auth, config.Conf.Email, []string{email}, msg)
	if err != nil {
		log.Fatal(err)
	}
}
