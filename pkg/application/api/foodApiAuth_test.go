package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApiAuthBodyIsCorrectlyBuilt(t *testing.T) {

	foodApiAuth := NewFoodApiAuth(nil)
	requestBody := string(foodApiAuth.buildAuthRequestBody("mail", "password"))

	assert.Equal(t, "{\"device_type\":\"ANDROID\",\"email\":\"mail\",\"password\":\"password\"}", requestBody)
}
