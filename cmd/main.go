package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/marc/get-food-to-go/pkg/application"
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
const RESERVATIONS_FILE_NAME = "pkg/resources/reservations.txt"

func init() {
	infrastructure.InitLogger()

	err := godotenv.Load(".env")
	if err != nil {
		zap.L().Error("Error loading .env file: ", zap.Error(err))
	}

	if os.Getenv("DB_USER") == "" {
		repository = persistance.NewFileRepository(BEARER_FILE_NAME, STORES_FILE_NAME, REFRESH_TOKEN_FILE_NAME, RESERVATIONS_FILE_NAME)
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
			storesString := strings.Join(application.StoresToString(availableStores), ", ")
			zap.L().Info("Found shop(s): " + storesString)

			stores := foodApi.SaveStores(availableStores)
			notificationService.SendAvailableStoresNotification(stores)

			reservations := storeService.GetReservations()
			if len(reservations) > 0 {
				foodApi.ReserveFood(stores, reservations)
			}
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
