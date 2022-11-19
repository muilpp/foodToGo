package main

import (
	"os"

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
const COUNTRIES_FILE_NAME = "pkg/resources/countryTestFile.txt"

func init() {
	infrastructure.InitLogger()

	err := godotenv.Load(".env")
	if err != nil {
		zap.L().Error("Error loading .env file: ", zap.Error(err))
	}

	if os.Getenv("DB_USER") == "" {
		repository = persistance.NewFileRepository(BEARER_FILE_NAME, STORES_FILE_NAME, REFRESH_TOKEN_FILE_NAME, COUNTRIES_FILE_NAME)
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
		var countryLessShops []domain.Store

		if len(availableStores) > 0 {
			countries := storeService.GetCountries()
			for _, country := range countries {
				stores := foodApi.FilterStoresByCountry(country.GetName(), availableStores)
				notificationService.SendNotification(stores, "_"+country.GetName())
				countryLessShops = append(countryLessShops, application.RemoveStoresFromSlice(availableStores, stores)...)
			}

			if len(countries) == 0 {
				countryLessShops = availableStores
			}
		}

		if len(countryLessShops) > 0 {
			notificationService.SendNotification(countryLessShops, "")
		}
	} else if executionType == "printGraph" {
		graphService = infrastructure.NewGraphService(repository)

		for _, country := range storeService.GetCountries() {
			graphService.PrintAllMonthlyReports(country.GetName())
			notificationService.SendTelegramMonthlyReports(country.GetName())
		}
	} else if executionType == "printGraphYear" {
		graphService = infrastructure.NewGraphService(repository)

		for _, country := range storeService.GetCountries() {
			graphService.PrintAllYearlyReports(country.GetName())
			notificationService.SendTelegramYearReports(country.GetName())
		}
	} else {
		zap.L().Warn("Wrong argument received in main function ", zap.String("Argument: ", executionType))
	}
}
