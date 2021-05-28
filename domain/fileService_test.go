package domain

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const storeTestFile = "storeTestFile.txt"
const bearerTestFile = "bearerTestFile.txt"

func TestMain(m *testing.M) {
	exitVal := m.Run()
	cleanup(bearerTestFile)
	cleanup(storeTestFile)
	os.Exit(exitVal)
}

func TestBearerIsCorrectlyReadFromFile(t *testing.T) {

	fs := NewFileService()

	fs.WriteBearerToFile(bearerTestFile, "AABBCC")
	bearer := fs.ReadBearerFromFile(bearerTestFile)

	assert.Equal(t, "AABBCC", bearer)
}

func TestStoresAreCorrectlyReadFromFile(t *testing.T) {

	fs := NewFileService()

	fs.WriteStoresToFile(storeTestFile, []string{"store1", "store2"})
	stores := fs.ReadStoresFromFile(storeTestFile)

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
