package ports

import "github.com/marc/get-food-to-go/pkg/domain"

type FoodService interface {
	GetStoresWithFood() []domain.Store
}

type FoodServiceAuth interface {
	Login() string
}

type FoodServiceToken interface {
	RefreshToken() string
}
