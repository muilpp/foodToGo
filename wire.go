package main

import (
	"github.com/google/wire"
	"github.com/marc/get-food-to-go/domain"
	"github.com/marc/get-food-to-go/domain/api"
)

func IntializeServices() (domain.FileService, domain.NotificationService, api.FoodApi) {
	wire.Build(domain.NewFileService())
	wire.Build(domain.NewNotificationService())
	wire.Build(api.NewFoodApi)

	return domain.FileService{}, domain.NotificationService{}, api.FoodApi{}
}
