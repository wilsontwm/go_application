package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"net/smtp"
	"strconv"
	"github.com/joho/godotenv"
)

type Request struct {
	From 	string
	To 		[]string
	Subject string
	Body 	string
}

type MailConfig struct {
	Server   string
	Port     int
	Email    string
	Password string
}

const (
	MIME = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

var config = MailConfig{}

// Initialize
func init() {
	err := godotenv.Load() //Load .env file
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	config.Server = os.Getenv("mail_server")
	config.Port, _ = strconv.Atoi(os.Getenv("mail_port"))
	config.Email = os.Getenv("mail_email")
	config.Password = os.Getenv("mail_password")
}

// Set the request
func NewRequest(to []string, subject string) *Request {
	return &Request{
		To: to,
		Subject: subject,
	}
}

// Parse the template for the email
func (r *Request) ParseTemplate(fileName string, data interface{}) error {
	t, err := template.ParseFiles(fileName)
	if err != nil {
		return err
	}

	buffer := new(bytes.Buffer)
	if err = t.Execute(buffer, data); err != nil {
		return err
	}

	r.Body = buffer.String()
	return nil
}

// Send email via SMTP
func (r *Request) SendMail() bool {
	body := "To: " + r.To[0] + "\r\nSubject: " + r.Subject + "\r\n" + MIME + "\r\n" + r.Body

	// Get mail server configuration
	SMTP := fmt.Sprintf("%s:%d", config.Server, config.Port)
	if err := smtp.SendMail(SMTP, smtp.PlainAuth("", config.Email, config.Password, config.Server), config.Email, r.To, []byte(body)); err != nil {
		log.Println(err)
		return false
	}

	return true
}

// Triggering point to send email
func (r *Request) Send(templateName string, data interface{}) {
	err := r.ParseTemplate(templateName, data)

	if err != nil {
		log.Println(err)
	}

	if ok := r.SendMail(); ok {
		log.Printf("Email has been successfully sent out to %s\n", r.To)
	} else {
		log.Printf("Email has failed to be sent out to %s\n", r.To)
	}
}