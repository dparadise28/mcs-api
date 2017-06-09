package tools

import (
	"crypto/tls"
	"encoding/json"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"models"
	"os"
	"sync"
	"time"
)

var (
	PlatformEmail = "yourfriends@mycorner.store"
	EmailPassword = ""
	queuePaths    = []string{
		"queues/email/PENDING/",
		"queues/email/RETRYING/",
		"queues/email/FAILED/",
		"queues/email/SUCCEEDED/",
	}
	fileMutex = &sync.Mutex{}
)

func EmailStructToGomailMsg(email *models.Email) *gomail.Message {
	m := gomail.NewMessage()
	m.SetHeader("From", PlatformEmail)
	m.SetHeader("To", email.To)
	m.SetHeader("Subject", email.Subject)
	m.SetBody("text/html", email.Body)
	return m
}

func SendEmail(email *models.Email) {
	m := EmailStructToGomailMsg(email)
	d := gomail.NewDialer("smtp.gmail.com", 587, PlatformEmail, EmailPassword)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
		log.Printf("email could not be sent")
	}
	//time.Sleep(time.Second * 4)
	log.Println("email sent")
}

func emailDaemon() {
	for email := range EmailQueue {
		mail, _ := json.Marshal(email)
		path := queuePaths[0] + time.Now().String()
		//fileMutex.Lock()
		f, err := os.Create(path)
		//log.Println("starting")
		if err != nil {
			log.Println(err)
			//fileMutex.Unlock()
			return
		}
		if _, err := f.Write(mail); err != nil {
			f.Close()
			log.Println(err)
			//fileMutex.Unlock()
			return
		}
		f.Close()
		//log.Println("closed file")
		//fileMutex.Unlock()
		go func() {
			log.Println("logging")
			Emailch <- &path
		}()
		//SendEmail(email)
	}
}

func emailerDaemon() {
	for email := range Emailch {
		emailStruct := models.Email{}
		b, err := ioutil.ReadFile(*email) // just pass the file name
		if err != nil {
			log.Println(err)
			return
		}
		if err := json.Unmarshal(b, &emailStruct); err != nil {
			log.Println(err)
			return
		}
		go SendEmail(&emailStruct)
	}
}

// buffered chanel for sending emails
var Emailch = make(chan *string, 50)

// unbuffered queue to send (unblocks the rest of the system
// while being blocked by the email sender daemon)
var EmailQueue = make(chan *models.Email, 1000000)

func initQueue() {
	for _, path := range queuePaths {
		os.MkdirAll(path, os.ModePerm)
	}
}

func init() {
	initQueue()
	go emailDaemon()
	go emailerDaemon()
}
