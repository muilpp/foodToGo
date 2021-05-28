package main

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
	"github.com/marc/get-food-to-go/domain"
	"github.com/marc/get-food-to-go/domain/api"
	"github.com/marc/get-food-to-go/resources"
)

var fileService domain.FileService
var notificationService domain.NotificationService
var foodApi api.FoodApi

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fileService = domain.NewFileService()
	notificationService = domain.NewNotificationService()
	foodApi = api.NewFoodApi(fileService)
}

func main() {
	availableStores := foodApi.GetStoresWithFood()

	if len(availableStores) > 0 {
		fileService.WriteStoresToFile(resources.StoresFileName, availableStores)

		storesString := strings.Join(availableStores, ", ")
		notificationService.SendMail(storesString)
		notificationService.SendTelegramMessage(storesString)
	}
}
