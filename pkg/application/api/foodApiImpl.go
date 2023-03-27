package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/marc/get-food-to-go/pkg/application"
	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"github.com/marc/get-food-to-go/pkg/infrastructure"
	"go.uber.org/zap"
)

const MAX_TRIES = 1
const NUM_ITEMS_TO_BOOK = 2

var currentTries int

type FoodApiImpl struct {
	foodAuth     ports.FoodServiceAuth
	foodToken    ports.FoodServiceToken
	storeService ports.StoreService
	userId       string
	latitude     string
	longitude    string
}

func NewFoodApi(authService ports.FoodServiceAuth, tokenService ports.FoodServiceToken, storeService ports.StoreService, userId string, latitude string, longitude string) FoodApiImpl {
	return FoodApiImpl{
		foodAuth:     authService,
		foodToken:    tokenService,
		storeService: storeService,
		userId:       userId,
		latitude:     latitude,
		longitude:    longitude,
	}
}

func (foodApi FoodApiImpl) GetStoresWithFood() []domain.Store {

	requestBody := foodApi.buildRequestBody()
	resp := foodApi.requestFood(requestBody, foodApi.getBearerToken())

	defer resp.Body.Close()

	if resp.StatusCode == 401 && currentTries < MAX_TRIES {
		currentTries++
		zap.L().Info("Unauthorized request, get new bearer")
		foodApi.foodToken.RefreshToken()
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
		"origin": {
			"latitude":` + foodApi.latitude + `,
			"longitude":` + foodApi.longitude + `
		},
		"radius": 5.0,
		"favorites_only": "true"
	}`)
}

func (foodApi FoodApiImpl) requestFood(requestBody []byte, bearerToken string) *http.Response {
	url := "https://apptoogoodtogo.com/api/item/v7/"
	method := "POST"

	payload := strings.NewReader(`{
	  "user_id": "` + foodApi.userId + `",
	  "origin": {"latitude": ` + foodApi.latitude + `, "longitude": ` + foodApi.longitude + `},
	  "radius": 5.0,
	  "favorites_only": "true"
  }`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "TooGoodToGo/22.10.0 (4665) (iPhone/iPhone XS Max; iOS 16.0.2; Scale/3.00/iOS)")
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "datadome=mDevyaYDHRxV1yK2HnPi_55SaGyz61eJ3v-Avcf2BT96tJlMDm0xKE~WRnSRpVCyzlgWkjt7esTjjp2ONe5mGe3HOcRcLr.sZYx7OKrBdZ7s5aOMWox.5hHZNZTSOuz")

	res, err := client.Do(req)
	if err != nil {
		zap.L().Fatal("Error requesting stores: ", zap.Error(err))
	}

	return res
}

func (foodApi FoodApiImpl) parseResponse(responseBody io.ReadCloser) domain.FoodJson {
	body, err := ioutil.ReadAll(responseBody)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	var responseStruct domain.FoodJson
	json.Unmarshal(body, &responseStruct)

	return responseStruct
}

func (foodApi FoodApiImpl) checkStoresInResponse(response domain.FoodJson) []domain.Store {
	var stores []domain.Store

	storesInFile := foodApi.storeService.GetStores()

	for _, item := range response.Items {
		if item.ItemsAvailable > 0 && item.Store.StoreName != "" {
			storeName := item.Store.StoreName

			if item.Item.Name != "" {
				storeName += " - " + item.Item.Name
			}

			if !application.StoresContainStoreName(storesInFile, storeName) {
				store := domain.NewStore(storeName, item.Item.ItemID, item.ItemsAvailable, time.Now())
				stores = append(stores, *store)
			}
		}
	}

	return stores
}

func (foodApi FoodApiImpl) SaveStores(availableStores []domain.Store) []domain.Store {
	var stores []domain.Store

	stores = append(stores, availableStores...)

	if len(stores) > 0 {
		foodApi.AddStores(stores)
		return stores
	}

	return []domain.Store{}
}

func (foodApi FoodApiImpl) AddStores(stores []domain.Store) {
	foodApi.storeService.AddStores(stores)
}

func (foodApi FoodApiImpl) ReserveFood(stores []domain.Store, reservations []string) {
	for _, reservation := range reservations {
		for _, store := range stores {
			if strings.Contains(strings.ToLower(store.GetName()), strings.ToLower(reservation)) {
				if store.GetItemsAvailable() < NUM_ITEMS_TO_BOOK {
					continue
				}

				resp := foodApi.requestFoodReservation(foodApi.getBearerToken(), store.GetItem(), NUM_ITEMS_TO_BOOK)

				defer resp.Body.Close()

				if resp.StatusCode != 200 {
					zap.L().Error("Bad response while trying to reserve food: ", zap.Int("status: ", resp.StatusCode), zap.Any("body", resp.Body))
				} else {
					reservationResponse := foodApi.parseReservationResponse(resp.Body)
					if reservationResponse.State == "SUCCESS" && reservationResponse.Order.State == "RESERVED" {
						notificationService := infrastructure.NewNotificationService()
						notificationService.SendReservationNotification(store.GetName(), reservationResponse.Order.ID)
					} else {
						zap.L().Error("Error at reservation, ", zap.String("State: ", reservationResponse.State), zap.String("Order state", reservationResponse.Order.State))
					}
				}
			}
		}
	}
}

func (foodApi FoodApiImpl) requestFoodReservation(bearerToken string, itemID string, numItems int) *http.Response {
	url := "https://apptoogoodtogo.com/api/order/v7/create/" + itemID
	method := "POST"

	payload := strings.NewReader(`{"item_count": ` + strconv.Itoa(numItems) + `}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "TooGoodToGo/22.10.0 (4665) (iPhone/iPhone XS Max; iOS 16.0.2; Scale/3.00/iOS)")
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Cookie", "datadome=mDevyaYDHRxV1yK2HnPi_55SaGyz61eJ3v-Avcf2BT96tJlMDm0xKE~WRnSRpVCyzlgWkjt7esTjjp2ONe5mGe3HOcRcLr.sZYx7OKrBdZ7s5aOMWox.5hHZNZTSOuz")

	res, err := client.Do(req)
	if err != nil {
		zap.L().Fatal("Error requesting stores: ", zap.Error(err))
	}

	return res
}

func (foodApi FoodApiImpl) parseReservationResponse(responseBody io.ReadCloser) domain.ReservationJson {
	body, err := ioutil.ReadAll(responseBody)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	var responseStruct domain.ReservationJson
	json.Unmarshal(body, &responseStruct)

	return responseStruct
}

func (foodApi FoodApiImpl) getBearerToken() string {
	bearerToken := foodApi.storeService.GetBearer()

	if bearerToken == "" {
		zap.L().Info("Current bearer empty, getting a new one")
		bearerToken = foodApi.foodToken.RefreshToken()
		zap.L().Info("New bearer: " + bearerToken)
	}

	return bearerToken
}
