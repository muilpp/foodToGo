package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	gomail "gopkg.in/mail.v2"
)

func sendMail(message string) {
	m := gomail.NewMessage()

	mailFrom := os.Getenv("MAIL_FROM")
	mailTo := os.Getenv("MAIL_TO")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSubject := os.Getenv("MAIL_SUBJECT")
	sliceTo := strings.Split(mailTo, ",")

	m.SetHeader("From", mailFrom)
	m.SetHeader("To", sliceTo...)
	m.SetHeader("Subject", mailSubject)
	m.SetBody("text/plain", message)

	d := gomail.NewDialer("smtp.gmail.com", 587, mailFrom, mailPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	if err := d.DialAndSend(m); err != nil {
		fmt.Println(err)
		panic(err)
	}

	return
}

func sendTelegramMessage(message string) {
	telegramToken := os.Getenv("TELEGRAM_API_TOKEN")
	telegramChatId := os.Getenv("TELEGRAM_CHAT_ID")

	requestUrl := "https://api.telegram.org/bot" + telegramToken + "/sendMessage?chat_id=" + telegramChatId + "&text=" + message

	_, err := http.Get(requestUrl)
	if err != nil {
		log.Fatalln(err)
	}
}
