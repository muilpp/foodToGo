package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiAuthBodyIsCorrectlyBuilt(t *testing.T) {

	foodApiAuth := NewFoodApiAuth(nil)
	requestBody := string(foodApiAuth.buildAuthRequestBody("mail"))

	assert.Equal(t, "{\"device_type\":\"ANDROID\",\"email\":\"mail\"}", requestBody)
}
