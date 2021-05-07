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

	bearerToken := readBearerFromFile()
	availableStores := getStoresWithFood(bearerToken)

	if len(availableStores) > 0 {
		writeStoresToFile(availableStores)

		storesString := strings.Join(availableStores, ",")
		sendMail(storesString)
		sendTelegramMessage(storesString)
	}
}
