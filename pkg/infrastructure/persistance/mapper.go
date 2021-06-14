package persistance

import (
	"github.com/marc/get-food-to-go/pkg/domain"
)

func StoreTableToStoreObject(storeTable StoreTable) *domain.Store {
	return domain.NewStore(storeTable.Store, "", 0)
}

func StoreObjectToStoreTable(store domain.Store) *StoreTable {
	return NewStoreTable(store.GetName())
}

func StoreTablesToStoreObjects(storeTables []StoreTable) []domain.Store {

	var stores []domain.Store

	for _, v := range storeTables {
		stores = append(stores, *domain.NewStore(v.Store, "", 0))
	}

	return stores
}

func StoreObjectsToStoreTables(storeObjects []domain.Store) []StoreTable {

	var stores []StoreTable

	for _, v := range storeObjects {
		stores = append(stores, *NewStoreTable(v.GetName()))
	}

	return stores
}
