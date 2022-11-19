package ports

import "github.com/marc/get-food-to-go/pkg/domain"

type FoodService interface {
	GetStoresWithFood() []domain.Store
	FilterStoresByCountry(countryCode string, availableStores []domain.Store) []domain.Store
	AddStores(stores []domain.Store)
}

type FoodServiceAuth interface {
	Login() string
}

type FoodServiceToken interface {
	RefreshToken() string
}
