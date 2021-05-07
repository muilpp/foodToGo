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
	availableStores := getFood(bearerToken)

	writeAvailableStoresToFile(availableStores)

	if len(availableStores) > 0 {
		sendMail(strings.Join(availableStores, ","))
		sendTelegramMessage(strings.Join(availableStores, ","))
	}
}
