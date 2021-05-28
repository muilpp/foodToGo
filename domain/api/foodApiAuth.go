package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/marc/get-food-to-go/domain"
	"github.com/marc/get-food-to-go/resources"
)

type FoodApiAuth struct {
	fileService domain.FileService
}

func NewFoodApiAuth(fs domain.FileService) *FoodApiAuth {
	return &FoodApiAuth{fs}
}

type AuthResponse struct {
	Token string `json:"access_token"`
}

func (apiAuth FoodApiAuth) GetAuthBearer() string {
	mail := os.Getenv("API_USER")
	password := os.Getenv("API_PASSWORD")

	values := map[string]string{"device_type": "ANDROID", "email": mail, "password": password}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post("https://apptoogoodtogo.com/api/auth/v2/loginByEmail", "application/json",
		bytes.NewBuffer(json_data))

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var authResponseObject AuthResponse
	json.Unmarshal(body, &authResponseObject)

	apiAuth.fileService.WriteBearerToFile(resources.BearerFileName, authResponseObject.Token)
	return authResponseObject.Token
}
