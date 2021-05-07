package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
)

const storesFileName = "availableStores.txt"
const bearerFileName = "authBearer.txt"

func readBearerFromFile() string {
	return readFile(bearerFileName)
}

func writeBearerToFile(bearer string) {
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

func readStoresFromFile() string {
	return readFile(storesFileName)
}

func writeStoresToFile(stores []string) {
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

func readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	return string(bytes.TrimSpace(content))
}
