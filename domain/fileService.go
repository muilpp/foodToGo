package domain

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

const storesFileName = "resources/availableStores.txt"
const bearerFileName = "resources/authBearer.txt"

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (fs *FileService) ReadBearerFromFile() string {
	return fs.readFile(bearerFileName)
}

func (fs *FileService) WriteBearerToFile(bearer string) {
	f, err := os.Create(bearerFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(bearer)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func (fs *FileService) ReadStoresFromFile() string {
	return fs.readFile(storesFileName)
}

func (fs *FileService) WriteStoresToFile(stores []string) {
	f, err := os.OpenFile(storesFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, store := range stores {
		_, err2 := f.WriteString(store + "\n")

		if err2 != nil {
			log.Fatal(err2)
		}
	}
}

func (fs *FileService) readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	return string(bytes.TrimSpace(content))
}
