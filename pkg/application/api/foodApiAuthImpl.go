package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"go.uber.org/zap"
)

type FoodApiAuthImpl struct {
	storeService ports.StoreService
}

func NewFoodApiAuth(storeService ports.StoreService) FoodApiAuthImpl {
	return FoodApiAuthImpl{storeService}
}

type AuthResponse struct {
	Token string `json:"access_token"`
}

func (apiAuth FoodApiAuthImpl) GetAuthBearer() string {
	json_data := apiAuth.buildAuthRequestBody(os.Getenv("API_USER"), os.Getenv("API_PASSWORD"))
	resp, err := http.Post("https://apptoogoodtogo.com/api/auth/v2/loginByEmail", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		zap.L().Error("Error getting bearer", zap.Error(err))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var authResponseObject AuthResponse
	json.Unmarshal(body, &authResponseObject)

	apiAuth.storeService.UpdateBearer(authResponseObject.Token)
	return authResponseObject.Token
}

func (apiAuth FoodApiAuthImpl) buildAuthRequestBody(mail string, password string) []byte {
	values := map[string]string{"device_type": "ANDROID", "email": mail, "password": password}
	json_data, err := json.Marshal(values)

	if err != nil {
		zap.L().Fatal("Could not marshall auth body", zap.Error(err))
	}

	return json_data
}
