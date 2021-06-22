package domain

type Store struct {
	name           string
	country        string
	itemsAvailable int
}

func NewStore(name string, country string, itemsAvailable int) *Store {
	return &Store{name, country, itemsAvailable}
}

func (s Store) GetName() string {
	return s.name
}

func (s Store) GetCountry() string {
	return s.country
}

func (s Store) GetItemsAvailable() int {
	return s.itemsAvailable
}

func (s Store) String() string {
	return s.GetName() + " (" + s.GetCountry() + ")"
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
