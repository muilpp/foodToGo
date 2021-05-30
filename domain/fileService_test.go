package domain

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const storeTestFile = "storeTestFile.txt"
const bearerTestFile = "bearerTestFile.txt"

var fileService PersistorService

func TestMain(m *testing.M) {
	fileService = NewFilePersistorService(bearerTestFile, storeTestFile)
	exitVal := m.Run()
	cleanup(bearerTestFile)
	cleanup(storeTestFile)
	os.Exit(exitVal)
}

func TestBearerIsCorrectlyReadFromFile(t *testing.T) {
	fileService.WriteBearer("AABBCC")
	bearer := fileService.ReadBearer()

	assert.Equal(t, "AABBCC", bearer)
}

func TestStoresAreCorrectlyReadFromFile(t *testing.T) {
	fileService.WriteStores([]string{"store1", "store2"})
	stores := fileService.ReadStores()

	assert.Equal(t, string("store1\nstore2"), stores)
}

func cleanup(fileToRemove string) {
	if _, err := os.Stat(fileToRemove); err == nil {
		e := os.Remove(fileToRemove)
		if e != nil {
			log.Fatal(e)
		}
	}
}
