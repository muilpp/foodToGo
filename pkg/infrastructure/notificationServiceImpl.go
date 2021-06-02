package infrastructure

import (
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"strings"

	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

type NotificationServiceImpl struct {
}

func NewNotificationService() *NotificationServiceImpl {
	return &NotificationServiceImpl{}
}

func (ns NotificationServiceImpl) SendMail(message string) {
	m := gomail.NewMessage()

	mailFrom := os.Getenv("MAIL_FROM")
	mailTo := os.Getenv("MAIL_TO")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSubject := os.Getenv("MAIL_SUBJECT")

	if mailFrom == "" || mailTo == "" || mailPassword == "" {
		return
	}

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
		zap.L().Error("Could not send mail", zap.Error(err))
	}
}

func (ns NotificationServiceImpl) SendTelegramMessage(message string) {
	telegramToken := os.Getenv("TELEGRAM_API_TOKEN")
	telegramChatId := os.Getenv("TELEGRAM_CHAT_ID")

	if telegramToken == "" || telegramChatId == "" {
		return
	}

	requestUrl := "https://api.telegram.org/bot" + telegramToken + "/sendMessage?chat_id=" + telegramChatId + "&text=" + message

	_, err := http.Get(requestUrl)
	if err != nil {
		log.Fatalln(err)
	}
}
