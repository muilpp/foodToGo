package api

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestCheckStoresInResponse(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"item\":{\"name\":\"cárnico\"},\"store\":{\"store_id\":\"47816s\",\"store_name\":\"Fruiteria\"},\"items_available\":0}, {\"store\": {\"store_name\": \"Carnisseria\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Pastisseria\"},\"items_available\": 1}]}}]}")
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
	bearerFile             string
	storeFile              string
	ReadStoresFromFileMock func() string
}

func newFileServiceMock(bearerFile string, storeFile string) fileServiceMock {
	return fileServiceMock{
		bearerFile:             bearerFile,
		storeFile:              storeFile,
		ReadStoresFromFileMock: func() string { return "Meat shop" },
	}
}

func (fs fileServiceMock) ReadBearer() string {
	return ""
}

func (fs fileServiceMock) ReadStores() string {
	return fs.ReadStoresFromFileMock()
}

func (fs fileServiceMock) WriteBearer(bearer string) {}

func (fs fileServiceMock) WriteStores(stores []string) {
}

func TestStoresNotAddedIfAlreadyPresentInFile(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	fs := newFileServiceMock("", "")
	foodApi := NewFoodApi(NewFoodApiAuth(fs), fs, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 2, len(stores), "2 shops expected, but only "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Bakery", stores[0], "Bakery expected, but found "+stores[0])
	assert.Equal(t, "Fish shop", stores[1], "Fish shop expected, but found "+stores[1])
}

func TestAllStoresAddedIfNoStoresInFile(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	fs := newFileServiceMock("", "")
	fs.ReadStoresFromFileMock = func() string {
		return ""
	}

	foodApi := NewFoodApi(NewFoodApiAuth(fs), fs, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 3, len(stores), "3 shops expected, but only "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Meat shop", stores[0], "Meat shop expected, but found "+stores[0])
	assert.Equal(t, "Bakery", stores[1], "Bakery expected, but found "+stores[1])
	assert.Equal(t, "Fish shop", stores[2], "Fish shop expected, but found "+stores[2])
}

func TestOnlyStoresWithItemsAvailableAdded(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 0}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 0}]}}]}")
	responseStruct := parseJsonResponse(response)

	fs := newFileServiceMock("", "")
	fs.ReadStoresFromFileMock = func() string {
		return ""
	}

	foodApi := NewFoodApi(NewFoodApiAuth(fs), fs, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 1, len(stores), "1 shop expected, but "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Meat shop", stores[0], "Meat shop expected, but found "+stores[0])
}

func TestStoresNameContainsItemName(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"item\":{\"name\":\"Meat\"},\"store\":{\"store_id\":\"47816s\",\"store_name\":\"Shop\"},\"items_available\":1}, {\"item\":{\"name\":\"Bread\"},\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 2}, {\"item\":{\"name\":\"Fish\"},\"store\": {\"store_name\": \"Shop\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	fs := newFileServiceMock("", "")
	fs.ReadStoresFromFileMock = func() string {
		return ""
	}

	foodApi := NewFoodApi(NewFoodApiAuth(fs), fs, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 3, len(stores), "3 shop expected, but "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Shop - Meat", stores[0], "Shop - Meat expected, but found "+stores[0])
	assert.Equal(t, "Bakery - Bread", stores[1], "Bakery - Bread expected, but found "+stores[1])
	assert.Equal(t, "Shop - Fish", stores[2], "Shop - Fish expected, but found "+stores[2])
}

func TestBuildRequestIsCorrectlyBuilt(t *testing.T) {
	fs := newFileServiceMock("", "")
	foodApi := NewFoodApi(NewFoodApiAuth(fs), fs, "123", "10", "20")
	bodyRequest := foodApi.buildRequestBody()

	assert.Equal(t, "{\n\t\t\"user_id\": \"123\",\n\t\t\"bucket_identifiers\": [\"Favorites\"],\n\t\t\"origin\": {\n\t\t\t\"latitude\":10,\n\t\t\t\"longitude\":20\n\t\t},\n\t\t\"radius\": 5.0,\n\t\t\"discover_experiments\": [\"WEIGHTED_ITEMS\"]\n\t}", string(bodyRequest), "Body request expected: "+""+", but found: ", string(bodyRequest))
}

type foodAuthMock struct {
	fileService domain.PersistorService
}

func newFoodAuthMock(fileService domain.PersistorService) *foodAuthMock {
	return &foodAuthMock{fileService}
}

func (authMock foodAuthMock) GetAuthBearer() string {
	return "bearerMock"
}

func TestNoStoresFoundWhenAuthBearerIsIncorrect(t *testing.T) {

	fileServiceMock := newFileServiceMock("", "")
	foodApi := NewFoodApi(newFoodAuthMock(fileServiceMock), fileServiceMock, "userId", "latitude", "longitude")
	stores := foodApi.GetStoresWithFood()

	assert.Equal(t, 0, len(stores))
}
