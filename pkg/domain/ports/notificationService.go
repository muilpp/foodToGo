package ports

import "github.com/marc/get-food-to-go/pkg/domain"

type NotificationService interface {
	SendAvailableStoresNotification(stores []domain.Store)
	SendReservationNotification(storeName string, itemID string)
	SendTelegramMonthlyReports()
	SendTelegramYearReports()
}
