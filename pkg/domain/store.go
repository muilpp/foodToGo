package domain

import "time"

type Store struct {
	name           string
	item           string
	itemsAvailable int
	createdAt      time.Time
}

func NewStore(name string, item string, itemsAvailable int, createdAt time.Time) *Store {
	return &Store{name, item, itemsAvailable, createdAt}
}

func (s Store) GetName() string {
	return s.name
}

func (s Store) GetItem() string {
	return s.item
}

func (s Store) GetItemsAvailable() int {
	return s.itemsAvailable
}

func (s Store) GetCreatedAt() time.Time {
	return s.createdAt
}

type StoreCounter struct {
	name  string
	total int
}

func NewStoreCounter(name string, total int) *StoreCounter {
	return &StoreCounter{name, total}
}

func (s StoreCounter) GetName() string {
	return s.name
}

func (s StoreCounter) GetTotal() int {
	return s.total
}

type ReservationStore struct {
	name          string
	alwaysReserve bool
}

func NewReservationStore(name string, alwaysReserve bool) *ReservationStore {
	return &ReservationStore{name, alwaysReserve}
}

func (s ReservationStore) GetName() string {
	return s.name
}

func (s ReservationStore) IsAlwaysReserve() bool {
	return s.alwaysReserve
}
