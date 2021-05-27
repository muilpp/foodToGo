package domain

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

type FileService interface {
	ReadBearerFromFile() string
	WriteBearerToFile(bearer string)
	ReadStoresFromFile() string
	WriteStoresToFile(stores []string)
}

type FileServiceImpl struct {
	StoresFileName string
	BearerFileName string
}

func NewFileService(bearerFileName string, storesFileName string) *FileServiceImpl {
	return &FileServiceImpl{
		BearerFileName: bearerFileName,
		StoresFileName: storesFileName,
	}
}

func (fs *FileServiceImpl) ReadBearerFromFile() string {
	return fs.readFile(fs.BearerFileName)
}

func (fs *FileServiceImpl) WriteBearerToFile(bearer string) {
	f, err := os.Create(fs.BearerFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(bearer)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func (fs *FileServiceImpl) ReadStoresFromFile() string {
	return fs.readFile(fs.StoresFileName)
}

func (fs *FileServiceImpl) WriteStoresToFile(stores []string) {
	f, err := os.OpenFile(fs.StoresFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

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
