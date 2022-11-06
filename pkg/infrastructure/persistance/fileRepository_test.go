package persistance

import (
	"os"
	"testing"

	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const storeTestFile = "storeTestFile.txt"
const bearerTestFile = "bearerTestFile.txt"
const refreshTokenTestFile = "refreshTokenTestFile.txt"

var fileService ports.StoreService

func TestMain(m *testing.M) {
	fileService = NewFileRepository(bearerTestFile, storeTestFile, refreshTokenTestFile)
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
	fileService.AddStores([]domain.Store{*domain.NewStore("store1", "", 0), *domain.NewStore("store2", "", 0)})
	stores := fileService.GetStores()

	zap.L().Info("Stores found: ", zap.Any("Stores: ", stores))
	assert.Equal(t, 2, len(stores))
	assert.Equal(t, "store1", stores[0].GetName())
	assert.Equal(t, "store2", stores[1].GetName())
}

func cleanup(fileToRemove string) {
	if _, err := os.Stat(fileToRemove); err == nil {
		e := os.Remove(fileToRemove)
		if e != nil {
			zap.L().Fatal(err.Error())
		}
	}
}
