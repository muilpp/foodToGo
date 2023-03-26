package ports

import (
	"strconv"
	"strings"

	"github.com/marc/get-food-to-go/pkg/domain"
	"go.uber.org/zap"
)

type StoreService interface {
	GetBearer() string
	UpdateBearer(bearer string)
	GetRefreshToken() string
	UpdateRefreshToken(token string)
	GetStores() []domain.Store
	AddStores(stores []domain.Store)
	GetReservations() []string
	ReserveFood([]domain.Store, []string)
}

type Repository interface {
	GetBearer() string
	UpdateBearer(bearer string)
	GetRefreshToken() string
	UpdateRefreshToken(token string)
	GetStores() []domain.Store
	GetStoresByTimesAppeared(frequency string) []domain.StoreCounter
	GetStoresByDayOfWeek(frequency string) []domain.StoreCounter
	GetStoresByHourOfDay(frequency string) []domain.StoreCounter
	AddStores(stores []domain.Store)
	GetReservations() []string
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

func (s storeService) GetRefreshToken() string {
	return s.repository.GetRefreshToken()
}

func (s storeService) UpdateRefreshToken(token string) {
	s.repository.UpdateRefreshToken(token)
}

func (s storeService) GetStores() []domain.Store {
	return s.repository.GetStores()
}

func (s storeService) AddStores(stores []domain.Store) {
	s.repository.AddStores(stores)
}

func (s storeService) GetReservations() []string {
	return s.repository.GetReservations()
}

func (s storeService) ReserveFood(stores []domain.Store, reservations []string) {
	zap.L().Info("Reservations: " + strconv.Itoa(len(reservations)))

	for _, reservation := range reservations {
		zap.L().Info("Reservation: " + reservation)
		for _, store := range stores {
			zap.L().Info("Store: " + store.GetName())

			if strings.Contains(strings.ToLower(store.GetName()), strings.ToLower(reservation)) {
				zap.L().Info("Tinc reserva per: " + store.GetName())
			}
		}
	}
}
