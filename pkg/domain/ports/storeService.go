package ports

import (
	"github.com/marc/get-food-to-go/pkg/domain"
)

type StoreService interface {
	GetBearer() string
	UpdateBearer(bearer string)
	GetStores() []domain.Store
	AddStores(stores []domain.Store)
}

type Repository interface {
	GetBearer() string
	UpdateBearer(bearer string)
	GetStores() []domain.Store
	GetStoresByTimesAppeared() []domain.StoreCounter
	GetStoresByDayOfWeek() []domain.StoreCounter
	GetStoresByHourOfDay() []domain.StoreCounter
	AddStores(stores []domain.Store)
}

type storeService struct {
	repository Repository
}

func NewStoreService(repository Repository) *storeService {
	return &storeService{repository}
}

func (s storeService) GetBearer() string {
	return s.repository.GetBearer()
}

func (s storeService) UpdateBearer(bearer string) {
	s.repository.UpdateBearer(bearer)
}

func (s storeService) GetStores() []domain.Store {
	return s.repository.GetStores()
}

func (s storeService) AddStores(stores []domain.Store) {
	s.repository.AddStores(stores)
}
