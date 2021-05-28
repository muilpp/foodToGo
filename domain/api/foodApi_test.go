package api

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckStoresInResponse(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Fruiteria\"},\"items_available\": 0}, {\"store\": {\"store_name\": \"Carnisseria\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Pastisseria\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	assert.Equal(t, 1, len(responseStruct.Groupings))
	assert.Equal(t, 3, len(responseStruct.Groupings[0].DiscoverBucket.Items), "Wrong number of shops found")
	assert.Equal(t, "Fruiteria", responseStruct.Groupings[0].DiscoverBucket.Items[0].Store.StoreName)
	assert.Equal(t, 0, responseStruct.Groupings[0].DiscoverBucket.Items[0].ItemsAvailable, "Fruiteria should have 0 items")
	assert.Equal(t, "Carnisseria", responseStruct.Groupings[0].DiscoverBucket.Items[1].Store.StoreName)
	assert.Equal(t, 2, responseStruct.Groupings[0].DiscoverBucket.Items[1].ItemsAvailable, "Carnisseria should have 2 items")
	assert.Equal(t, "Pastisseria", responseStruct.Groupings[0].DiscoverBucket.Items[2].Store.StoreName)
	assert.Equal(t, 1, responseStruct.Groupings[0].DiscoverBucket.Items[2].ItemsAvailable, "Pastisseria should have 1 item")
}

func parseJsonResponse(response []byte) FoodJson {
	var responseStruct FoodJson
	err := json.Unmarshal(response, &responseStruct)

	if err != nil {
		panic("Could not unmarshal json response: " + err.Error())
	}

	return responseStruct
}

type fileServiceMock struct {
	ReadStoresFromFileMock func(bearerFile string) string
}

func newFileServiceMock() fileServiceMock {
	return fileServiceMock{
		ReadStoresFromFileMock: func(bearerFile string) string { return "Meat shop" },
	}
}

func (fs fileServiceMock) ReadBearerFromFile(bearerFile string) string {
	return ""
}

func (fs fileServiceMock) ReadStoresFromFile(bearerFile string) string {
	return fs.ReadStoresFromFileMock(bearerFile)
}

func (fs fileServiceMock) WriteBearerToFile(bearerFile string, bearer string) {}

func (fs fileServiceMock) WriteStoresToFile(bearerFile string, stores []string) {
}

func TestStoresNotAddedIfAlreadyPresentInFile(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	fs := newFileServiceMock()

	foodApi := NewFoodApi(fs)

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 2, len(stores), "2 shops expected, but only "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Bakery", stores[0], "Bakery expected, but found "+stores[0])
	assert.Equal(t, "Fish shop", stores[1], "Fish shop expected, but found "+stores[1])
}

func TestAllStoresAddedIfNoStoresInFile(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	fs := newFileServiceMock()
	fs.ReadStoresFromFileMock = func(bearerFile string) string {
		return ""
	}

	foodApi := NewFoodApi(fs)

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 3, len(stores), "3 shops expected, but only "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Meat shop", stores[0], "Meat shop expected, but found "+stores[0])
	assert.Equal(t, "Bakery", stores[1], "Bakery expected, but found "+stores[1])
	assert.Equal(t, "Fish shop", stores[2], "Fish shop expected, but found "+stores[2])
}

func TestOnlyStoresWithItemsAvailableSaved(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 0}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 0}]}}]}")
	responseStruct := parseJsonResponse(response)

	fs := newFileServiceMock()
	fs.ReadStoresFromFileMock = func(bearerFile string) string {
		return ""
	}

	foodApi := NewFoodApi(fs)

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 1, len(stores), "1 shop expected, but "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Meat shop", stores[0], "Meat shop expected, but found "+stores[0])
}
