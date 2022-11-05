package persistance

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/marc/get-food-to-go/pkg/domain"
)

type FileRepository struct {
	bearerFileName       string
	storesFileName       string
	refreshTokenFileName string
}

func NewFileRepository(bearerFileName string, storesFileName string, refreshTokenFileName string) *FileRepository {
	return &FileRepository{bearerFileName: bearerFileName, storesFileName: storesFileName, refreshTokenFileName: refreshTokenFileName}
}

func (fs *FileRepository) GetBearer() string {
	return fs.readFile(fs.bearerFileName)
}

func (fs *FileRepository) UpdateBearer(bearer string) {
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

func (fs *FileRepository) GetRefreshToken() string {
	return fs.readFile(fs.refreshTokenFileName)
}

func (fs *FileRepository) UpdateRefreshToken(token string) {
	f, err := os.Create(fs.refreshTokenFileName)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(token)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func (fs *FileRepository) GetStores() []domain.Store {
	var stores []domain.Store
	storeStrings := strings.Split(fs.readFile(fs.storesFileName), "\n")

	for _, storeString := range storeStrings {
		stores = append(stores, *domain.NewStore(storeString, "", 0))
	}

	return stores
}

func (fs *FileRepository) GetStoresByTimesAppeared(frequency string) []domain.StoreCounter {
	return []domain.StoreCounter{}
}

func (fs *FileRepository) GetStoresByDayOfWeek(frequency string) []domain.StoreCounter {
	return []domain.StoreCounter{}
}

func (fs *FileRepository) GetStoresByHourOfDay(frequency string) []domain.StoreCounter {
	return []domain.StoreCounter{}
}

func (fs *FileRepository) AddStores(stores []domain.Store) {
	f, err := os.OpenFile(fs.storesFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for _, store := range stores {
		_, err2 := f.WriteString(store.GetName() + "\n")

		if err2 != nil {
			log.Fatal(err2)
		}
	}
}

func (fs *FileRepository) readFile(fileName string) string {
	content, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}

	return string(bytes.TrimSpace(content))
}
