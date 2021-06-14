package api

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/marc/get-food-to-go/pkg/application"
	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"go.uber.org/zap"
)

const MAX_TRIES = 1

var currentTries int

type FoodApiImpl struct {
	foodAuth          ports.FoodServiceAuth
	repositoryService ports.Repository
	userId            string
	latitude          string
	longitude         string
}

func NewFoodApi(authService ports.FoodServiceAuth, fs ports.Repository, userId string, latitude string, longitude string) FoodApiImpl {
	return FoodApiImpl{
		foodAuth:          authService,
		repositoryService: fs,
		userId:            userId,
		latitude:          latitude,
		longitude:         longitude,
	}
}

func (foodApi FoodApiImpl) GetStoresWithFood() []domain.Store {

	bearerToken := foodApi.repositoryService.GetBearer()

	if bearerToken == "" {
		bearerToken = foodApi.foodAuth.GetAuthBearer()
		zap.L().Info("Current bearer empty, getting a new one")
	}

	requestBody := foodApi.buildRequestBody()
	resp := foodApi.requestFood(requestBody, bearerToken)

	defer resp.Body.Close()

	if resp.StatusCode == 401 && currentTries < MAX_TRIES {
		currentTries++
		zap.L().Info("Unauthorized request, get new bearer")
		foodApi.foodAuth.GetAuthBearer()
		return foodApi.GetStoresWithFood()
	} else if resp.StatusCode != 200 {
		zap.L().Error("Bad response while getting stores: ", zap.Int("status: ", resp.StatusCode), zap.Any("body", resp.Body))
		return []domain.Store{}
	}

	response := foodApi.parseResponse(resp.Body)
	return foodApi.checkStoresInResponse(response)
}

func (foodApi FoodApiImpl) buildRequestBody() []byte {

	return []byte(`{
		"user_id": "` + foodApi.userId + `",
		"bucket_identifiers": ["Favorites"],
		"origin": {
			"latitude":` + foodApi.latitude + `,
			"longitude":` + foodApi.longitude + `
		},
		"radius": 5.0,
		"discover_experiments": ["WEIGHTED_ITEMS"]
	}`)
}

func (foodApi FoodApiImpl) requestFood(requestBody []byte, bearerToken string) *http.Response {
	req, err := http.NewRequest("POST", "https://apptoogoodtogo.com/api/item/v7/discover", bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "PostmanRuntime/7.26.8")
	req.Header.Add("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		zap.L().Fatal("Error requesting stores: ", zap.Error(err))
	}

	return resp
}

func (foodApi FoodApiImpl) parseResponse(responseBody io.ReadCloser) domain.FoodJson {
	body, err := ioutil.ReadAll(responseBody)
	if err != nil {
		log.Fatalln(err)
	}

	var responseStruct domain.FoodJson
	json.Unmarshal(body, &responseStruct)

	return responseStruct
}

func (foodApi FoodApiImpl) checkStoresInResponse(response domain.FoodJson) []domain.Store {
	var stores []domain.Store

	storesInFile := foodApi.repositoryService.GetStores()

	for _, grouping := range response.Groupings {
		for _, item := range grouping.DiscoverBucket.Items {
			if item.ItemsAvailable > 0 {
				storeName := item.Store.StoreName
				if item.Item.Name != "" {
					storeName += " - " + item.Item.Name
				}

				if !application.StoresContainStoreName(storesInFile, storeName) {
					store := domain.NewStore(storeName, item.Store.StoreLocation.Address.Country.IsoCode, item.ItemsAvailable)
					stores = append(stores, *store)
				}
			}
		}
	}

	return stores
}
