package main

import (
	"fmt"
	"log"
	"net/smtp"
)

func sendEmail(msg Message) {
	smtpServer := "smtp-mail.outlook.com:587"
	auth := smtp.PlainAuth("", "220392@astanait.edu.kz", "Sarzhan123", "smtp-mail.outlook.com")

	subject := fmt.Sprintf("%s message", msg.Type)
	body := fmt.Sprintf("Subject: %s\n\n%s", subject, msg.Message)

	hostname := "smtp-mail.outlook.com"
	auth = LoginAuth(hostname, "220392@astanait.edu.kz", "Sarzhan123", "smtp-mail.outlook.com")
	// Send the email
	err := smtp.SendMail(smtpServer, auth, "220392@astanait.edu.kz", []string{msg.Recipient}, []byte(body))
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
