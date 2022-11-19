package persistance

import (
	"time"

	"github.com/marc/get-food-to-go/pkg/domain"
)

func StoreTableToStoreObject(storeTable StoreTable) *domain.Store {
	return domain.NewStore(storeTable.Store, "", "", 0, time.Now())
}

func StoreObjectToStoreTable(store domain.Store) *StoreTable {
	return NewStoreTable(store.GetName(), store.GetCountry())
}

func StoreTablesToStoreObjects(storeTables []StoreTable) []domain.Store {

	var stores []domain.Store

	for _, v := range storeTables {
		stores = append(stores, *domain.NewStore(v.Store, "", "", 0, v.CreatedAt))
	}

	return stores
}

func StoreObjectsToStoreTables(storeObjects []domain.Store) []StoreTable {

	var stores []StoreTable

	for _, v := range storeObjects {
		stores = append(stores, *NewStoreTable(v.GetName(), v.GetCountry()))
	}

	return stores
}

func StoreTableCountResultsToStoreCounterObjects(storeTables []Result) []domain.StoreCounter {

	var storesCounter []domain.StoreCounter

	for _, v := range storeTables {
		storesCounter = append(storesCounter, *domain.NewStoreCounter(v.Element, v.Total))
	}

	return storesCounter
}

func CountryTableToCountryObject(countryTable []CountryTable) []domain.Country {

	var countries []domain.Country

	for _, v := range countryTable {
		countries = append(countries, *domain.NewCountry(v.Country))
	}

	return countries
}
