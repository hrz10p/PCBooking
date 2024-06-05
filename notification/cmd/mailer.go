package main

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func sendEmail(msg Message) {
	smtpServer := os.Getenv("SMTP_SERVER")
	smtpPort := os.Getenv("SMTP_PORT")
	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	auth := smtp.PlainAuth("", email, password, smtpServer)

	subject := fmt.Sprintf("%s message", msg.Type)
	body := fmt.Sprintf("Subject: %s\nMIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n%s", subject, msg.Message)

	hostname := smtpServer
	auth = LoginAuth(hostname, email, password, smtpServer)
	// Send the email
	err := smtp.SendMail(smtpServer+":"+smtpPort, auth, email, []string{msg.Recipient}, []byte(body))
	if err != nil {
		log.Printf("Error sending email to %s: %v", msg.Recipient, err)
	} else {
		log.Printf("Email sent to %s", msg.Recipient)
	}
}

func LoginAuth(hostname, username, password, host string) smtp.Auth {
	return &loginAuth{hostname, username, password, host}
}

type loginAuth struct {
	hostname, username, password, host string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, fmt.Errorf("unknown message: %s", fromServer)
		}
	}
	return nil, nil
}
