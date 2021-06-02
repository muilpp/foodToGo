package ports

type NotificationService interface {
	SendMail(message string)
	SendTelegramMessage(message string)
}
