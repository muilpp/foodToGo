package main

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/marc/get-food-to-go/domain"
	"github.com/marc/get-food-to-go/domain/api"
)

const storesFileName = "resources/availableStores.txt"
const bearerFileName = "resources/authBearer.txt"

var fileService domain.FileService
var notificationService domain.NotificationService
var foodApi api.FoodApi

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fileService = domain.NewFileService(bearerFileName, storesFileName)
	notificationService = domain.NewNotificationService()
	foodApi = api.NewFoodApi(fileService)

	bearerToken := fileService.ReadBearerFromFile()
	availableStores := foodApi.GetStoresWithFood(bearerToken)

	if len(availableStores) > 0 {
		fileService.WriteStoresToFile(availableStores)

		storesString := strings.Join(availableStores, ", ")
		notificationService.SendMail(storesString)
		notificationService.SendTelegramMessage(storesString)
	}
}
