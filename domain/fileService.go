package domain

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

type PersistorService interface {
	ReadBearer() string
	WriteBearer(bearer string)
	ReadStores() string
	WriteStores(stores []string)
}

type FilePersistor struct {
	bearerFileName string
	storesFileName string
}

func NewFilePersistorService(bearerFileName string, storesFileName string) *FilePersistor {
	return &FilePersistor{bearerFileName: bearerFileName, storesFileName: storesFileName}
}

func (fs *FilePersistor) ReadBearer() string {
	return fs.readFile(fs.bearerFileName)
}

func (fs *FilePersistor) WriteBearer(bearer string) {
	f, err := os.Create(fs.bearerFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(bearer)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func (fs *FilePersistor) ReadStores() string {
	return fs.readFile(fs.storesFileName)
}

func (fs *FilePersistor) WriteStores(stores []string) {
	f, err := os.OpenFile(fs.storesFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

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

func (fs *FilePersistor) readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	return string(bytes.TrimSpace(content))
}
