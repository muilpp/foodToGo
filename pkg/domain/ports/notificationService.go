package ports

import "github.com/marc/get-food-to-go/pkg/domain"

type NotificationService interface {
	SendNotification(stores []domain.Store, telegramToken string, telegramChatId int64)
	SendTelegramMonthlyReports(countryCode string, telegramToken string, telegramChatId int64)
	SendTelegramYearReports(countryCode string, telegramToken string, telegramChatId int64)
}
