package domain

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

type FileService interface {
	ReadBearerFromFile(bearerFile string) string
	WriteBearerToFile(bearerFile string, bearer string)
	ReadStoresFromFile(bearerFile string) string
	WriteStoresToFile(bearerFile string, stores []string)
}

type FileServiceImpl struct {
}

func NewFileService() *FileServiceImpl {
	return &FileServiceImpl{}
}

func (fs *FileServiceImpl) ReadBearerFromFile(bearerFile string) string {
	return fs.readFile(bearerFile)
}

func (fs *FileServiceImpl) WriteBearerToFile(bearerFile string, bearer string) {
	f, err := os.Create(bearerFile)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(bearer)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func (fs *FileServiceImpl) ReadStoresFromFile(storeFile string) string {
	return fs.readFile(storeFile)
}

func (fs *FileServiceImpl) WriteStoresToFile(storeFile string, stores []string) {
	f, err := os.OpenFile(storeFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

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

func (fs *FileServiceImpl) readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	return string(bytes.TrimSpace(content))
}
