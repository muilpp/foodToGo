package application

import (
	"time"

	"github.com/marc/get-food-to-go/pkg/domain"
)

func StoreToString(store domain.Store) string {
	return store.GetName()
}

func StringToStore(name string, country string, itemsAvailable int, createdAt time.Time) *domain.Store {
	return domain.NewStore(name, country, "", itemsAvailable, createdAt)
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
		stores = append(stores, *domain.NewStore(store, "", "", 0, time.Now()))
	}

	return stores
}

func StoresContainStoreName(stores []domain.Store, storeName string) bool {

	for _, store := range stores {
		if store.GetName() == storeName && time.Since(store.GetCreatedAt()).Minutes() < 120 {
			return true
		}
	}

	return false
}

func RemoveStoresFromSlice(originalStores []domain.Store, storesToRemove []domain.Store) []domain.Store {
	var updatedStores []domain.Store

	for _, os := range originalStores {
		removeStore := false
		for _, sr := range storesToRemove {

			if os.GetName() == sr.GetName() {
				removeStore = true
			}
		}

		if !removeStore {
			updatedStores = append(updatedStores, os)
		}
	}

	return updatedStores
}
