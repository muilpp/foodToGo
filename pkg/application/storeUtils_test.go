package application

import (
	"testing"
	"time"

	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/stretchr/testify/assert"
)

func TestStoreIsConvertedToString(t *testing.T) {
	store := domain.NewStore("Meat shop", "ES", 1, time.Now())
	stringStore := StoreToString(*store)

	assert.Equal(t, "Meat shop", stringStore)
}

func TestStoreSliceIsConvertedToStringSlice(t *testing.T) {

	stores := []domain.Store{*domain.NewStore("Meat shop", "ES", 1, time.Now()), *domain.NewStore("Fish shop", "FR", 2, time.Now())}
	stringStores := StoresToString(stores)

	assert.Equal(t, 2, len(stringStores))
	assert.Equal(t, "Meat shop", stringStores[0])
	assert.Equal(t, "Fish shop", stringStores[1])
}

func TestStringIsConvertedToStore(t *testing.T) {
	store := StringToStore("Meat shop", "ES", 1, time.Now())

	assert.Equal(t, "Meat shop", store.GetName())
	assert.Equal(t, "ES", store.GetCountry())
	assert.Equal(t, 1, store.GetItemsAvailable())
}

func TestStringSliceIsConvertedToStoreSlice(t *testing.T) {
	storeString := []string{"Meat shop", "Fish shop"}
	stores := StringsToStores(storeString)

	assert.Equal(t, 2, len(stores))
	assert.Equal(t, "Meat shop", stores[0].GetName())
	assert.Equal(t, "Fish shop", stores[1].GetName())
}

func TestStoreSliceContainsStoreName(t *testing.T) {
	stores := []domain.Store{*domain.NewStore("Meat shop", "ES", 1, time.Now()), *domain.NewStore("Fish shop", "FR", 2, time.Now()), *domain.NewStore("Candy shop", "IT", 3, time.Now())}

	assert.True(t, StoresContainStoreName(stores, "Meat shop"))
	assert.True(t, StoresContainStoreName(stores, "Fish shop"))
	assert.True(t, StoresContainStoreName(stores, "Candy shop"))
	assert.False(t, StoresContainStoreName(stores, "Random shop"))
}
