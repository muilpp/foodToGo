package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/marc/get-food-to-go/pkg/domain"
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
	PollingId string `json:"polling_id"`
}

func (apiAuth FoodApiAuthImpl) Login() string {
	json_data := apiAuth.buildAuthRequestBody(os.Getenv("API_USER"), os.Getenv("API_PASSWORD"))
	req, err := http.NewRequest("POST", "https://apptoogoodtogo.com/api/auth/v3/authByEmail/", bytes.NewBuffer(json_data))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+apiAuth.storeService.GetBearer())
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept-Language", "en-GB")
	req.Header.Add("User-Agent", "TooGoodToGo/22.10.0 (4665) (iPhone/iPhone XS Max; iOS 16.0.2; Scale/3.00/iOS)")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")

	client := &http.Client{}
	resp, err := client.Do(req)

	if resp.StatusCode > 400 {
		zap.L().Info("Auth response status: " + strconv.Itoa(resp.StatusCode))
		zap.L().Info("Auth response message: " + resp.Status)
		panic("Authentication error, let's stop here...")
	}

	if err != nil {
		zap.L().Error("Error getting bearer", zap.Error(err))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	var authResponseObject AuthResponse
	json.Unmarshal(body, &authResponseObject)

	zap.L().Info("Polling id: " + authResponseObject.PollingId)

	apiAuth.storeService.UpdateBearer(authResponseObject.PollingId)
	return authResponseObject.PollingId
}

func (apiAuth FoodApiAuthImpl) buildAuthRequestBody(mail string, password string) []byte {
	values := map[string]string{"device_type": "ANDROID", "email": mail}
	json_data, err := json.Marshal(values)

	if err != nil {
		zap.L().Fatal("Could not marshall auth body", zap.Error(err))
	}

	return json_data
}

func (apiAuth FoodApiAuthImpl) RefreshToken() string {
	url := "https://apptoogoodtogo.com/api/auth/v3/token/refresh"
	method := "POST"

	payload := strings.NewReader(`{"refresh_token":"` + apiAuth.storeService.GetRefreshToken() + `"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("User-Agent", "TooGoodToGo/22.10.0 (4665) (iPhone/iPhone XS Max; iOS 16.0.2; Scale/3.00/iOS)")
	req.Header.Add("Accept-Language", "en-GB")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+apiAuth.storeService.GetBearer())
	req.Header.Add("Cookie", "datadome=mMvQpdhREIj7mhipR-3Z4ArhebKomdbRdCR-y8p9966vTyZGQCaZv-fJJa~.pAcqGgoiDYqb4lxYmY__LMzfug2NIj02rOll~RxyWFlqSUPgpNsUw2r78Mk.yfk_LmM")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	zap.L().Info("Bearer response status: " + res.Status)

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		zap.L().Fatal(err.Error())
	}

	var tokenResponse domain.RefreshTokenResponse
	json.Unmarshal(body, &tokenResponse)

	if tokenResponse.AccessToken != "" {
		apiAuth.storeService.UpdateBearer(tokenResponse.AccessToken)
	}

	if tokenResponse.RefreshToken != "" {
		apiAuth.storeService.UpdateRefreshToken(tokenResponse.RefreshToken)
	}

	zap.L().Info("Bearer TTL" + strconv.Itoa(tokenResponse.AccessTokenTTLSeconds))

	return string(tokenResponse.AccessToken)
}
