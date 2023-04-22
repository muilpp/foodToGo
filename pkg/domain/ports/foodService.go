package ports

import "github.com/marc/get-food-to-go/pkg/domain"

type FoodService interface {
	GetStoresWithFood() []domain.Store
	SaveStores(availableStores []domain.Store) []domain.Store
	AddStores(stores []domain.Store)
	ReserveFood([]domain.Store, []domain.ReservationStore)
}

type FoodServiceAuth interface {
	Login() string
}

type FoodServiceToken interface {
	RefreshToken() string
}
