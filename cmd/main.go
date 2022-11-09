package main

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/marc/get-food-to-go/pkg/application"
	"github.com/marc/get-food-to-go/pkg/application/api"
	"github.com/marc/get-food-to-go/pkg/domain"
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

			spanishStores := filterShops("ES", availableStores)
			germanStores := filterShops("DE", availableStores)

			telegramSpanishToken := os.Getenv("TELEGRAM_API_TOKEN")
			telegramSpanishChatId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID"), 10, 64)
			processShops(spanishStores, telegramSpanishToken, telegramSpanishChatId, true)

			telegramGermanToken := os.Getenv("TELEGRAM_API_TOKEN_GERMAN")
			telegramGermanChatId, _ := strconv.ParseInt(os.Getenv("TELEGRAM_CHAT_ID_GERMAN"), 10, 64)
			processShops(germanStores, telegramGermanToken, telegramGermanChatId, false)
		}
	} else if executionType == "printGraph" {
		graphService = infrastructure.NewGraphService(repository)
		graphService.PrintAllMonthlyReports()
		notificationService.SendTelegramMonthlyReports()
	} else if executionType == "printGraphYear" {
		graphService = infrastructure.NewGraphService(repository)
		graphService.PrintAllYearlyReports()
		notificationService.SendTelegramYearReports()

	} else {
		zap.L().Warn("Wrong argument received in main function ", zap.String("Argument: ", executionType))
	}
}

func filterShops(countryCode string, availableStores []domain.Store) []domain.Store {
	var stores []domain.Store

	for _, store := range availableStores {
		if store.GetCountry() == countryCode {
			stores = append(stores, store)
		}
	}

	return stores
}

func processShops(stores []domain.Store, telegramToken string, telegramChatId int64, addToDatabase bool) {
	if len(stores) > 0 {
		storesString := strings.Join(application.StoresToString(stores), ", ")
		zap.L().Info("Found shop(s): " + storesString)

		if addToDatabase {
			storeService.AddStores(stores)
		}

		notificationService.SendMail(storesString)
		notificationService.SendTelegramMessage(storesString, telegramToken, telegramChatId)
	}
}
