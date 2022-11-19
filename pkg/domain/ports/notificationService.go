package ports

import "github.com/marc/get-food-to-go/pkg/domain"

type NotificationService interface {
	SendNotification(stores []domain.Store, countryName string)
	SendTelegramMonthlyReports(countryCode string)
	SendTelegramYearReports(countryCode string)
}
