package api

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"github.com/stretchr/testify/assert"
)

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
		bearerFile: bearerFile,
		storeFile:  storeFile,
		ReadStoresFromFileMock: func() []domain.Store {
			return []domain.Store{*domain.NewStore("Osakii - Mainz - Abendbuffet", "DE", "", 1, time.Now())}
		},
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

func (fs RepositoryMock) GetCountries() []domain.Country {
	return []domain.Country{}
}

func TestStoresResponseParse(t *testing.T) {
	response := []byte("{\"items\":[{\"item\":{\"item_id\":\"796374\",\"name\":\"Abendbuffet\"},\"store\":{\"store_id\":\"768562\",\"store_name\":\"Villa Wow - Berlin\"},\"display_name\":\"Villa Wow - Berlin\",\"items_available\":3},{\"item\":{\"item_id\":\"796374\",\"name\":\"Entrepans\"},\"store\":{\"store_id\":\"768562\",\"store_name\":\"Mini Golf\"},\"display_name\":\"Mini Golf\",\"items_available\":2}]}")
	responseStruct := parseJsonResponse(response)

	repoMock := newRepositoryMock("", "")
	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 2, len(stores), "2 shops expected, but only "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Villa Wow - Berlin - Abendbuffet", stores[0].GetName(), "Villa Wow - Berlin expected, but found "+stores[0].GetName())
	assert.Equal(t, 3, stores[0].GetItemsAvailable())
	assert.Equal(t, "Mini Golf - Entrepans", stores[1].GetName(), "Mini Golf expected, but found "+stores[1].GetName())
	assert.Equal(t, 2, stores[1].GetItemsAvailable())
}

func TestOnlyStoresWithItemsAvailableAdded(t *testing.T) {
	response := []byte("{\"items\":[{\"item\":{\"item_id\":\"796374\",\"name\":\"Abendbuffet\"},\"store\":{\"store_id\":\"768562\",\"store_name\":\"Osakii - Mainz\"},\"display_name\":\"Osakii - Mainz\",\"items_available\":3},{\"item\":{\"item_id\":\"796374\",\"name\":\"Entrepans\"},\"store\":{\"store_id\":\"768562\",\"store_name\":\"Mini Golf\"},\"display_name\":\"Mini Golf\",\"items_available\":0}]}")
	responseStruct := parseJsonResponse(response)

	repoMock := newRepositoryMock("", "")
	repoMock.ReadStoresFromFileMock = func() []domain.Store {
		return []domain.Store{}
	}

	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 1, len(stores), "1 shop expected, but "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Osakii - Mainz - Abendbuffet", stores[0].GetName(), "Osakii - Mainz - Abendbuffet expected, but found "+stores[0].GetName())
	assert.Equal(t, 3, stores[0].GetItemsAvailable())
}

func TestStoresNotAddedIfAlreadyPresentInFile(t *testing.T) {
	response := []byte("{\"items\":[{\"item\":{\"item_id\":\"796374\",\"name\":\"Abendbuffet\"},\"store\":{\"store_id\":\"768562\",\"store_name\":\"Osakii - Mainz\"},\"display_name\":\"Osakii - Mainz\",\"items_available\":3},{\"item\":{\"item_id\":\"796374\",\"name\":\"Entrepans\"},\"store\":{\"store_id\":\"768562\",\"store_name\":\"Mini Golf\"},\"display_name\":\"Mini Golf\",\"items_available\":2},{\"item\":{\"item_id\":\"796374\",\"name\":\"Bocadillos\"},\"store\":{\"store_id\":\"768562\",\"store_name\":\"Frankfurt 92\"},\"display_name\":\"Frankfurt 92\",\"items_available\":3}]}")
	responseStruct := parseJsonResponse(response)

	repoMock := newRepositoryMock("", "")
	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "", "", "")

	stores := foodApi.checkStoresInResponse(responseStruct)

	assert.Equal(t, 2, len(stores), "2 shops expected, but only "+strconv.Itoa(len(stores))+" found")
	assert.Equal(t, "Mini Golf - Entrepans", stores[0].GetName(), "Mini Golf expected, but found "+stores[0].GetName())
	assert.Equal(t, 2, stores[0].GetItemsAvailable())
	assert.Equal(t, "Frankfurt 92 - Bocadillos", stores[1].GetName(), "Frankfurt 92 - Bocadillos expected, but found "+stores[1].GetName())
	assert.Equal(t, 3, stores[1].GetItemsAvailable())
}

func TestBuildRequestIsCorrectlyBuilt(t *testing.T) {
	repoMock := newRepositoryMock("", "")
	foodApi := NewFoodApi(NewFoodApiAuth(repoMock), NewFoodApiAuth(repoMock), repoMock, "123", "10", "20")
	bodyRequest := foodApi.buildRequestBody()

	assert.Equal(t, "{\n\t\t\"user_id\": \"123\",\n\t\t\"origin\": {\n\t\t\t\"latitude\":10,\n\t\t\t\"longitude\":20\n\t\t},\n\t\t\"radius\": 5.0,\n\t\t\"favorites_only\": \"true\"\n\t}", string(bodyRequest), "Body request expected: "+""+", but found: ", string(bodyRequest))
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
