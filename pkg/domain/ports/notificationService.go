package ports

import "github.com/marc/get-food-to-go/pkg/domain"

type NotificationService interface {
	SendNotification(stores []domain.Store, countryName string)
	SendTelegramMonthlyReportsDeclaredCountry(countryCode string)
	SendTelegramMonthlyReports(telegramCountryCode string, filenameCountryCode string)
	SendTelegramYearReportsDeclaredCountry(countryCode string)
	SendTelegramYearReports(telegramCountryCode string, filenameCountryCode string)
}
