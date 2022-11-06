package api

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"github.com/stretchr/testify/assert"
)

// func TestCheckStoresInResponse(t *testing.T) {
// 	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"item\":{\"name\":\"cárnico\"},\"store\":{\"store_id\":\"47816s\",\"store_name\":\"Fruiteria\"},\"items_available\":0}, {\"store\": {\"store_name\": \"Carnisseria\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Pastisseria\"},\"items_available\": 1}]}}]}")
// 	responseStruct := parseJsonResponse(response)

// 	assert.Equal(t, 1, len(responseStruct.Groupings))
// 	assert.Equal(t, 3, len(responseStruct.Groupings[0].DiscoverBucket.Items), "Wrong number of shops found")
// 	assert.Equal(t, "Fruiteria", responseStruct.Groupings[0].DiscoverBucket.Items[0].Store.StoreName)
// 	assert.Equal(t, 0, responseStruct.Groupings[0].DiscoverBucket.Items[0].ItemsAvailable, "Fruiteria should have 0 items")
// 	assert.Equal(t, "Carnisseria", responseStruct.Groupings[0].DiscoverBucket.Items[1].Store.StoreName)
// 	assert.Equal(t, 2, responseStruct.Groupings[0].DiscoverBucket.Items[1].ItemsAvailable, "Carnisseria should have 2 items")
// 	assert.Equal(t, "Pastisseria", responseStruct.Groupings[0].DiscoverBucket.Items[2].Store.StoreName)
// 	assert.Equal(t, 1, responseStruct.Groupings[0].DiscoverBucket.Items[2].ItemsAvailable, "Pastisseria should have 1 item")
// }

func parseJsonResponse(response []byte) domain.FoodJson {
	var responseStruct domain.FoodJson
	err := json.Unmarshal(response, &responseStruct)

	if err != nil {
		panic("Could not unmarshal json response: " + err.Error())
	}

	return responseStruct
}

type RepositoryMock struct {
	bearerFile             string
	storeFile              string
	ReadStoresFromFileMock func() []domain.Store
}

func newRepositoryMock(bearerFile string, storeFile string) RepositoryMock {
	return RepositoryMock{
		bearerFile:             bearerFile,
		storeFile:              storeFile,
		ReadStoresFromFileMock: func() []domain.Store { return []domain.Store{*domain.NewStore("Meat shop", "ES", 1)} },
	}
}

func (fs RepositoryMock) GetBearer() string {
	return ""
}

func (fs RepositoryMock) GetStores() []domain.Store {
	return fs.ReadStoresFromFileMock()
}

func (fs RepositoryMock) GetStoresByTimesAppeared() []domain.Store {
	return []domain.Store{}
}

func (fs RepositoryMock) GetStoresByDayOfWeek() []domain.Store {
	return []domain.Store{}
}

func (fs RepositoryMock) GetStoresByHourOfDay() []domain.Store {
	return []domain.Store{}
}

func (fs RepositoryMock) UpdateBearer(bearer string) {}

func (fs RepositoryMock) GetRefreshToken() string {
	return ""
}
func (fs RepositoryMock) UpdateRefreshToken(bearer string) {}

func (fs RepositoryMock) AddStores(stores []domain.Store) {
}

func TestStoresNotAddedIfAlreadyPresentInFile(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	repoMock := newRepositoryMock("", "")
	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 2, len(stores), "2 shops expected, but only "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Bakery", stores[0].GetName(), "Bakery expected, but found "+stores[0].GetName())
	assert.Equal(t, 2, stores[0].GetItemsAvailable())
	assert.Equal(t, "Fish shop", stores[1].GetName(), "Fish shop expected, but found "+stores[1].GetName())
	assert.Equal(t, 1, stores[1].GetItemsAvailable())
}

func TestAllStoresAddedIfNoStoresInFile(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	repoMock := newRepositoryMock("", "")
	repoMock.ReadStoresFromFileMock = func() []domain.Store {
		return []domain.Store{}
	}

	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 3, len(stores), "3 shops expected, but only "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Meat shop", stores[0].GetName(), "Meat shop expected, but found "+stores[0].GetName())
	assert.Equal(t, 1, stores[0].GetItemsAvailable())
	assert.Equal(t, "Bakery", stores[1].GetName(), "Bakery expected, but found "+stores[1].GetName())
	assert.Equal(t, 2, stores[1].GetItemsAvailable())
	assert.Equal(t, "Fish shop", stores[2].GetName(), "Fish shop expected, but found "+stores[2].GetName())
	assert.Equal(t, 1, stores[2].GetItemsAvailable())
}

func TestOnlyStoresWithItemsAvailableAdded(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Meat shop\"},\"items_available\": 1}, {\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 0}, {\"store\": {\"store_name\": \"Fish shop\"},\"items_available\": 0}]}}]}")
	responseStruct := parseJsonResponse(response)

	repoMock := newRepositoryMock("", "")
	repoMock.ReadStoresFromFileMock = func() []domain.Store {
		return []domain.Store{}
	}

	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 1, len(stores), "1 shop expected, but "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Meat shop", stores[0].GetName(), "Meat shop expected, but found "+stores[0].GetName())
	assert.Equal(t, 1, stores[0].GetItemsAvailable())
}

func TestStoresNameContainsItemName(t *testing.T) {
	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"item\":{\"name\":\"Meat\"},\"store\":{\"store_id\":\"47816s\",\"store_name\":\"Shop\"},\"items_available\":1}, {\"item\":{\"name\":\"Bread\"},\"store\": {\"store_name\": \"Bakery\"},\"items_available\": 2}, {\"item\":{\"name\":\"Fish\"},\"store\": {\"store_name\": \"Shop\"},\"items_available\": 1}]}}]}")
	responseStruct := parseJsonResponse(response)

	repoMock := newRepositoryMock("", "")
	repoMock.ReadStoresFromFileMock = func() []domain.Store {
		return []domain.Store{}
	}

	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 3, len(stores), "3 shop expected, but "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Shop - Meat", stores[0].GetName(), "Shop - Meat expected, but found "+stores[0].GetName())
	assert.Equal(t, 1, stores[0].GetItemsAvailable())
	assert.Equal(t, "Bakery - Bread", stores[1].GetName(), "Bakery - Bread expected, but found "+stores[1].GetName())
	assert.Equal(t, 2, stores[1].GetItemsAvailable())
	assert.Equal(t, "Shop - Fish", stores[2].GetName(), "Shop - Fish expected, but found "+stores[2].GetName())
	assert.Equal(t, 1, stores[2].GetItemsAvailable())
}

func TestBuildRequestIsCorrectlyBuilt(t *testing.T) {
	repoMock := newRepositoryMock("", "")
	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "123", "10", "20")
	bodyRequest := foodApi.buildRequestBody()

	assert.Equal(t, "{\n\t\t\"user_id\": \"123\",\n\t\t\"bucket_identifiers\": [\"Favorites\"],\n\t\t\"origin\": {\n\t\t\t\"latitude\":10,\n\t\t\t\"longitude\":20\n\t\t},\n\t\t\"radius\": 5.0,\n\t\t\"discover_experiments\": [\"WEIGHTED_ITEMS\"]\n\t}", string(bodyRequest), "Body request expected: "+""+", but found: ", string(bodyRequest))
}

type foodAuthMock struct {
	storeService ports.StoreService
}

func newFoodAuthMock(storeService ports.StoreService) *foodAuthMock {
	return &foodAuthMock{storeService}
}

func (authMock foodAuthMock) Login() string {
	return "bearerMock"
}

func TestNoStoresFoundWhenAuthBearerIsIncorrect(t *testing.T) {

	repositoryMock := newRepositoryMock("", "")
	foodApi := NewFoodApi(newFoodAuthMock(repositoryMock), NewFoodApiAuth(repositoryMock), repositoryMock, "userId", "latitude", "longitude")
	stores := foodApi.GetStoresWithFood()

	assert.Equal(t, 0, len(stores))
}
