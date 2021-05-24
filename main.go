package main

import (
	"log"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	fileService, notificationService, foodApi := IntializeServices()

	bearerToken := fileService.ReadBearerFromFile()
	availableStores := foodApi.GetStoresWithFood(bearerToken)

	if len(availableStores) > 0 {
		fileService.WriteStoresToFile(availableStores)

		storesString := strings.Join(availableStores, ",")
		notificationService.SendMail(storesString)
		notificationService.SendTelegramMessage(storesString)
	}
}
