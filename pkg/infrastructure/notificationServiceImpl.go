package infrastructure

import (
	"crypto/tls"
	"os"
	"path"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

type NotificationServiceImpl struct {
}

func NewNotificationService() *NotificationServiceImpl {
	return &NotificationServiceImpl{}
}

func (ns NotificationServiceImpl) sendMail(message string) {
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

func (ns NotificationServiceImpl) sendTelegramMessage(message string, telegramToken string, telegramChatId int64) {

	if telegramToken == "" || telegramChatId == 0 {
		return
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	bot.Debug = true

	if err != nil {
		zap.L().Panic("Could not get telegram API instance", zap.Error(err))
	}

	msg := tgbotapi.NewMessage(telegramChatId, message)
	_, err2 := bot.Send(msg)

	if err2 != nil {
		zap.L().Error("Message not sent", zap.Error(err2))
	}
}

func (ns NotificationServiceImpl) SendTelegramMonthlyReports(countryCode string, telegramToken string, telegramChatId int64) {
	telegramChatId, bot := getTelegramCredentials(countryCode, telegramToken, telegramChatId)

	fileDir, _ := os.Getwd()
	ns.sendFile(countryCode+ports.FOOD_CHART_BY_STORE_MONTHLY, bot, fileDir, telegramChatId)
	ns.sendFile(countryCode+ports.FOOD_CHART_BY_DAY_OF_WEEK_MONTHLY, bot, fileDir, telegramChatId)
	ns.sendFile(countryCode+ports.FOOD_CHART_BY_HOUR_OF_DAY_MONTHLY, bot, fileDir, telegramChatId)
}

func (ns NotificationServiceImpl) SendTelegramYearReports(countryCode string, telegramToken string, telegramChatId int64) {
	telegramChatId, bot := getTelegramCredentials(countryCode, telegramToken, telegramChatId)

	fileDir, _ := os.Getwd()
	ns.sendFile(countryCode+ports.FOOD_CHART_BY_STORE_YEARLY, bot, fileDir, telegramChatId)
	ns.sendFile(countryCode+ports.FOOD_CHART_BY_DAY_OF_WEEK_YEARLY, bot, fileDir, telegramChatId)
	ns.sendFile(countryCode+ports.FOOD_CHART_BY_HOUR_OF_DAY_YEARLY, bot, fileDir, telegramChatId)
}

func getTelegramCredentials(countryCode string, telegramToken string, telegramChatId int64) (int64, *tgbotapi.BotAPI) {

	if telegramToken == "" || telegramChatId == 0 {
		zap.L().Panic("Got empty telegram credentials")
	}

	bot, err := tgbotapi.NewBotAPI(telegramToken)
	bot.Debug = true

	if err != nil {
		zap.L().Panic("Could not get telegram API instance", zap.Error(err))
	}

	return telegramChatId, bot
}

func (ns NotificationServiceImpl) sendFile(fileName string, bot *tgbotapi.BotAPI, fileDir string, chatId int64) {
	filePath := path.Join(fileDir, fileName)

	msg := tgbotapi.NewPhoto(chatId, filePath)
	_, err2 := bot.Send(msg)

	if err2 != nil {
		zap.L().Error("Document not sent", zap.Error(err2))
	}
}

func (ns NotificationServiceImpl) SendNotification(storesString string, telegramToken string, telegramChatId int64) {
	if len(storesString) > 0 {
		ns.sendMail(storesString)
		ns.sendTelegramMessage(storesString, telegramToken, telegramChatId)
	}
}
