package ports

type NotificationService interface {
	SendMail(message string)
	SendTelegramMessage(message string, telegramToken string, telegramChatId int64)
	SendTelegramMonthlyReports()
	SendTelegramYearReports()
}
