package tools

import (
	"crypto/tls"
	"gopkg.in/gomail.v2"
	"log"
)

var (
	PlatformEmail = "yourfriends@mycorner.store"
	EmailPassword = ""
)

func SendEmailValidation(to string, subject string, body string) (success bool) {
	m := gomail.NewMessage()
	m.SetHeader("From", PlatformEmail)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer("smtp.office365.com", 587, "yourfriends@mycorner.store", "!&GDrmfIX0yYv!Ec")
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		return false
	}
	return true
}
