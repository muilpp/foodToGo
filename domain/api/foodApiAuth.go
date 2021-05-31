package api

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/marc/get-food-to-go/domain"
)

type FoodApiAuth struct {
	fileService domain.PersistorService
}

func NewFoodApiAuth(fs domain.PersistorService) *FoodApiAuth {
	return &FoodApiAuth{fs}
}

type AuthResponse struct {
	Token string `json:"access_token"`
}

func (apiAuth FoodApiAuth) GetAuthBearer() string {
	json_data := apiAuth.buildAuthRequestBody(os.Getenv("API_USER"), os.Getenv("API_PASSWORD"))
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

	apiAuth.fileService.WriteBearer(authResponseObject.Token)
	return authResponseObject.Token
}

func (apiAuth FoodApiAuth) buildAuthRequestBody(mail string, password string) []byte {
	values := map[string]string{"device_type": "ANDROID", "email": mail, "password": password}
	json_data, err := json.Marshal(values)

	if err != nil {
		log.Fatal(err)
	}

	return json_data
}
