package ports

type NotificationService interface {
	SendNotification(storesString string, telegramToken string, telegramChatId int64)
	SendTelegramMonthlyReports(countryCode string, telegramToken string, telegramChatId int64)
	SendTelegramYearReports(countryCode string, telegramToken string, telegramChatId int64)
}
