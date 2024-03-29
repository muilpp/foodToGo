package infrastructure

import (
	"crypto/tls"
	"os"
	"path"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"go.uber.org/zap"
	gomail "gopkg.in/mail.v2"
)

type NotificationServiceImpl struct {
}

func NewNotificationService() *NotificationServiceImpl {
	return &NotificationServiceImpl{}
}

func (ns NotificationServiceImpl) sendMail(stores []domain.Store) {
	m := gomail.NewMessage()

	mailFrom := os.Getenv("MAIL_FROM")
	mailTo := os.Getenv("MAIL_TO")
	mailPassword := os.Getenv("MAIL_PASSWORD")
	mailSubject := os.Getenv("MAIL_SUBJECT")

	if mailFrom == "" || mailTo == "" || mailPassword == "" {
		return
	}

	var message string
	for _, s := range stores {
		message += s.GetName() + "\n"
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

func (ns NotificationServiceImpl) sendTelegramMessage(stores []domain.Store) {

	telegramChatId, bot := getTelegramCredentials()

	for _, store := range stores {
		link := "https://share.toogoodtogo.com/item/" + store.GetItem()
		msg := tgbotapi.NewMessage(telegramChatId, store.GetName()+"\n"+link, true)
		_, err2 := bot.Send(msg)

		if err2 != nil {
			zap.L().Error("Store message not sent", zap.Error(err2))
		}
	}
}

func (ns NotificationServiceImpl) SendTelegramMonthlyReports() {
	telegramChatId, bot := getTelegramCredentials()

	fileDir, _ := os.Getwd()
	ns.sendFile(ports.FOOD_CHART_BY_STORE_MONTHLY, bot, fileDir, telegramChatId)
	ns.sendFile(ports.FOOD_CHART_BY_DAY_OF_WEEK_MONTHLY, bot, fileDir, telegramChatId)
	ns.sendFile(ports.FOOD_CHART_BY_HOUR_OF_DAY_MONTHLY, bot, fileDir, telegramChatId)
}

func (ns NotificationServiceImpl) SendTelegramYearReports() {
	telegramChatId, bot := getTelegramCredentials()

	fileDir, _ := os.Getwd()
	ns.sendFile(ports.FOOD_CHART_BY_STORE_YEARLY, bot, fileDir, telegramChatId)
	ns.sendFile(ports.FOOD_CHART_BY_DAY_OF_WEEK_YEARLY, bot, fileDir, telegramChatId)
	ns.sendFile(ports.FOOD_CHART_BY_HOUR_OF_DAY_YEARLY, bot, fileDir, telegramChatId)
}

func getTelegramCredentials() (int64, *tgbotapi.BotAPI) {

	telegramToken := os.Getenv("TELEGRAM_API_TOKEN")
	telegramChatId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)

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
		zap.L().Error("Document "+fileName+" not sent", zap.Error(err2))
	}
}

func (ns NotificationServiceImpl) SendAvailableStoresNotification(stores []domain.Store) {
	if len(stores) > 0 {
		ns.sendMail(stores)
		ns.sendTelegramMessage(stores)
	}
}

func (ns NotificationServiceImpl) SendReservationNotification(storeName string, itemID string) {
	telegramChatId, bot := getTelegramCredentials()

	msg := tgbotapi.NewMessage(telegramChatId, "Reservation made for "+storeName+", item id: "+itemID, true)
	_, err2 := bot.Send(msg)

	if err2 != nil {
		zap.L().Error("Reservation message not sent", zap.Error(err2))
	}
}
