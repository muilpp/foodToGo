package api

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckStoresInResponse(t *testing.T) {

	response := []byte("{\"groupings\": [{\"discover_bucket\": {\"items\": [{\"store\": {\"store_name\": \"Fruiteria\"},\"items_available\": 0}, {\"store\": {\"store_name\": \"Carnisseria\"},\"items_available\": 2}, {\"store\": {\"store_name\": \"Pastisseria\"},\"items_available\": 1}]}}]}")

	var responseStruct FoodJson
	err := json.Unmarshal(response, &responseStruct)

	if err != nil {
		fmt.Println("Error: ", err)
	}

	assert.Equal(t, 1, len(responseStruct.Groupings))
	assert.Equal(t, 3, len(responseStruct.Groupings[0].DiscoverBucket.Items), "Wrong number of shops found")
	assert.Equal(t, "Fruiteria", responseStruct.Groupings[0].DiscoverBucket.Items[0].Store.StoreName)
	assert.Equal(t, 0, responseStruct.Groupings[0].DiscoverBucket.Items[0].ItemsAvailable, "Fruiteria should have 0 items")
	assert.Equal(t, "Carnisseria", responseStruct.Groupings[0].DiscoverBucket.Items[1].Store.StoreName)
	assert.Equal(t, 2, responseStruct.Groupings[0].DiscoverBucket.Items[1].ItemsAvailable, "Carnisseria should have 2 items")
	assert.Equal(t, "Pastisseria", responseStruct.Groupings[0].DiscoverBucket.Items[2].Store.StoreName)
	assert.Equal(t, 1, responseStruct.Groupings[0].DiscoverBucket.Items[2].ItemsAvailable, "Pastisseria should have 1 item")
}
