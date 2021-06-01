package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/api"
	"github.com/marc/get-food-to-go/pkg/infrastructure"
	"go.uber.org/zap"
)

var logger *zap.Logger
var fileService domain.PersistorService
var authService api.FoodApiAuth
var notificationService domain.NotificationService
var foodApi api.FoodApi

const STORES_FILE_NAME = "pkg/resources/availableStores.txt"
const BEARER_FILE_NAME = "pkg/resources/authBearer.txt"

func init() {
	logger := infrastructure.InitLogger()
	zap.ReplaceGlobals(logger)

	err := godotenv.Load(".env")
	if err != nil {
		logger.Fatal("Error loading .env file: ", zap.Error(err))
	}

	fileService = domain.NewFilePersistorService(BEARER_FILE_NAME, STORES_FILE_NAME)
	authService = api.NewFoodApiAuth(fileService)
	notificationService = domain.NewNotificationService()
	foodApi = api.NewFoodApi(authService, fileService, os.Getenv("APP_USER_ID"), os.Getenv("LATITUDE"), os.Getenv("LONGITUDE"))
}

func main() {
	availableStores := foodApi.GetStoresWithFood()

	if len(availableStores) > 0 {
		zap.L().Info("Found ", zap.Strings("stores: ", availableStores))
		fileService.WriteStores(availableStores)

		storesString := strings.Join(availableStores, ", ")
		notificationService.SendMail(storesString)
		notificationService.SendTelegramMessage(storesString)
	}
}
