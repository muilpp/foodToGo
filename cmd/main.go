package main

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/marc/get-food-to-go/pkg/application/api"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"github.com/marc/get-food-to-go/pkg/infrastructure"
	"github.com/marc/get-food-to-go/pkg/infrastructure/persistance"
	"go.uber.org/zap"
)

var storeService ports.StoreService
var graphService ports.GraphService
var repository ports.Repository

var authService ports.FoodServiceAuth
var tokenService ports.FoodServiceToken
var notificationService ports.NotificationService
var foodApi ports.FoodService

const STORES_FILE_NAME = "pkg/resources/availableStores.txt"
const BEARER_FILE_NAME = "pkg/resources/authBearer.txt"
const REFRESH_TOKEN_FILE_NAME = "pkg/resources/refreshToken.txt"

func init() {
	infrastructure.InitLogger()

	err := godotenv.Load(".env")
	if err != nil {
		zap.L().Error("Error loading .env file: ", zap.Error(err))
	}

	if os.Getenv("DB_USER") == "" {
		repository = persistance.NewFileRepository(BEARER_FILE_NAME, STORES_FILE_NAME, REFRESH_TOKEN_FILE_NAME)
	} else {
		repository = persistance.NewMysqlRepository(os.Getenv("DB_USER"), os.Getenv("DB_PWD"), os.Getenv("DB_IP"), os.Getenv("DB_NAME"))
	}

	storeService = ports.NewStoreService(repository)

	authService = api.NewFoodApiAuth(storeService)
	tokenService = api.NewFoodApiAuth(storeService)
	notificationService = infrastructure.NewNotificationService()
	foodApi = api.NewFoodApi(authService, tokenService, storeService, os.Getenv("APP_USER_ID"), os.Getenv("LATITUDE"), os.Getenv("LONGITUDE"))
}

func main() {
	executionType := os.Args[1]

	if executionType == "getFood" {
		availableStores := foodApi.GetStoresWithFood()

		if len(availableStores) > 0 {
			for _, country := range storeService.GetCountries() {
				stores := foodApi.FilterStoresByCountry(country.GetName(), availableStores)
				telegramToken, telegramChatId := getTelegramCredentials(country.GetName())
				notificationService.SendNotification(stores, telegramToken, telegramChatId)
			}
		}
	} else if executionType == "printGraph" {
		graphService = infrastructure.NewGraphService(repository)

		for _, country := range storeService.GetCountries() {
			graphService.PrintAllMonthlyReports(country.GetName())
			telegramToken, telegramChatId := getTelegramCredentials(country.GetName())
			notificationService.SendTelegramMonthlyReports(country.GetName(), telegramToken, telegramChatId)
		}
	} else if executionType == "printGraphYear" {
		graphService = infrastructure.NewGraphService(repository)

		for _, country := range storeService.GetCountries() {
			graphService.PrintAllYearlyReports(country.GetName())
			telegramToken, telegramChatId := getTelegramCredentials(country.GetName())
			notificationService.SendTelegramYearReports(country.GetName(), telegramToken, telegramChatId)
		}
	} else {
		zap.L().Warn("Wrong argument received in main function ", zap.String("Argument: ", executionType))
	}
}

func getTelegramCredentials(countryCode string) (string, int64) {
	telegramToken := os.Getenv("TELEGRAM_API_TOKEN_" + countryCode)
	telegramChatId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID_"+countryCode), 10, 64)

	return telegramToken, telegramChatId
}
