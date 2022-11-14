package persistance

import (
	"os"
	"testing"
	"time"

	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const storeTestFile = "storeTestFile.txt"
const bearerTestFile = "bearerTestFile.txt"
const refreshTokenTestFile = "refreshTokenTestFile.txt"
const countryTestFile = "countryTestFile.txt"

var fileService ports.StoreService

func TestMain(m *testing.M) {
	fileService = NewFileRepository(bearerTestFile, storeTestFile, refreshTokenTestFile, countryTestFile)
	exitVal := m.Run()
	cleanup(bearerTestFile)
	cleanup(storeTestFile)
	os.Exit(exitVal)
}

func TestBearerIsCorrectlyReadFromFile(t *testing.T) {
	fileService.UpdateBearer("AABBCC")
	bearer := fileService.GetBearer()

	assert.Equal(t, "AABBCC", bearer)
}

func TestStoresAreCorrectlyReadFromFile(t *testing.T) {
	fileService.AddStores([]domain.Store{*domain.NewStore("store1", "", "", 0, time.Now()), *domain.NewStore("store2", "", "", 0, time.Now())})
	stores := fileService.GetStores()

	zap.L().Info("Stores found: ", zap.Any("Stores: ", stores))
	assert.Equal(t, 2, len(stores))
	assert.Equal(t, "store1", stores[0].GetName())
	assert.Equal(t, "store2", stores[1].GetName())
}

func TestCountriesAreCorrectlyReadFromFile(t *testing.T) {
	countries := []domain.Country{*domain.NewCountry("ES"), *domain.NewCountry("DE"), *domain.NewCountry("FR")}

	writeCountriesToFile(countryTestFile, countries)
	countriesFound := fileService.GetCountries()

	assert.Equal(t, 3, len(countriesFound))
	assert.Equal(t, "ES", countriesFound[0].GetName())
	assert.Equal(t, "DE", countriesFound[1].GetName())
	assert.Equal(t, "FR", countriesFound[2].GetName())

	cleanup(countryTestFile)
}

func writeCountriesToFile(fileName string, countries []domain.Country) {
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		zap.L().Fatal(err.Error())
	}

	defer f.Close()

	for _, country := range countries {
		_, err2 := f.WriteString(country.GetName() + "\n")

		if err2 != nil {
			zap.L().Fatal(err2.Error())
		}
	}
}

func cleanup(fileToRemove string) {
	if _, err := os.Stat(fileToRemove); err == nil {
		e := os.Remove(fileToRemove)
		if e != nil {
			zap.L().Fatal(err.Error())
		}
	}
}
