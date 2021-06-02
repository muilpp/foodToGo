package application

import "github.com/marc/get-food-to-go/pkg/domain"

func StoreToString(store domain.Store) string {
	return store.GetName()
}

func StringToStore(name string, country string, itemsAvailable int) *domain.Store {
	return domain.NewStore(name, country, itemsAvailable)
}

func StoresToString(stores []domain.Store) []string {
	var storesSlice []string

	for _, store := range stores {
		storesSlice = append(storesSlice, store.GetName())
	}

	return storesSlice
}

func StringsToStores(stringStores []string) []domain.Store {
	var stores []domain.Store

	for _, store := range stringStores {
		stores = append(stores, *domain.NewStore(store, "", 0))
	}

	return stores
}

func StoresContainStoreName(stores []domain.Store, storeName string) bool {

	for _, store := range stores {
		if store.GetName() == storeName {
			return true
		}
	}

	return false
}
