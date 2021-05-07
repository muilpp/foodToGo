package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const MAX_TRIES = 1

var currentTries int

type FoodJson struct {
	ItemAvailabilityStatus string `json:"item_availability_status"`
	Groupings              []struct {
		Type           string `json:"type"`
		DiscoverBucket struct {
			BucketIdentifier string `json:"bucket_identifier"`
			Items            []struct {
				Item struct {
					ItemID string `json:"item_id"`
					Price  struct {
						Code       string `json:"code"`
						MinorUnits int    `json:"minor_units"`
						Decimals   int    `json:"decimals"`
					} `json:"price"`
					SalesTaxes   []interface{} `json:"sales_taxes"`
					CoverPicture struct {
						PictureID  string `json:"picture_id"`
						CurrentURL string `json:"current_url"`
					} `json:"cover_picture"`
					LogoPicture struct {
						PictureID  string `json:"picture_id"`
						CurrentURL string `json:"current_url"`
					} `json:"logo_picture"`
					Name                   string        `json:"name"`
					Description            string        `json:"description"`
					CanUserSupplyPackaging bool          `json:"can_user_supply_packaging"`
					PackagingOption        string        `json:"packaging_option"`
					CollectionInfo         string        `json:"collection_info"`
					DietCategories         []interface{} `json:"diet_categories"`
					ItemCategory           string        `json:"item_category"`
					Badges                 []struct {
						BadgeType   string `json:"badge_type"`
						RatingGroup string `json:"rating_group"`
						Percentage  int    `json:"percentage"`
						UserCount   int    `json:"user_count"`
						MonthCount  int    `json:"month_count"`
					} `json:"badges"`
					PositiveRatingReasons []string `json:"positive_rating_reasons"`
					AverageOverallRating  struct {
						AverageOverallRating float64 `json:"average_overall_rating"`
						RatingCount          int     `json:"rating_count"`
						MonthCount           int     `json:"month_count"`
					} `json:"average_overall_rating"`
					FavoriteCount int  `json:"favorite_count"`
					Buffet        bool `json:"buffet"`
				} `json:"item"`
				Store struct {
					StoreID       string `json:"store_id"`
					StoreName     string `json:"store_name"`
					Branch        string `json:"branch"`
					Description   string `json:"description"`
					TaxIdentifier string `json:"tax_identifier"`
					Website       string `json:"website"`
					StoreLocation struct {
						Address struct {
							Country struct {
								IsoCode string `json:"iso_code"`
								Name    string `json:"name"`
							} `json:"country"`
							AddressLine string `json:"address_line"`
							City        string `json:"city"`
							PostalCode  string `json:"postal_code"`
						} `json:"address"`
						Location struct {
							Longitude float64 `json:"longitude"`
							Latitude  float64 `json:"latitude"`
						} `json:"location"`
					} `json:"store_location"`
					LogoPicture struct {
						PictureID  string `json:"picture_id"`
						CurrentURL string `json:"current_url"`
					} `json:"logo_picture"`
					Distance     float64 `json:"distance"`
					CoverPicture struct {
						PictureID  string `json:"picture_id"`
						CurrentURL string `json:"current_url"`
					} `json:"cover_picture"`
				} `json:"store"`
				DisplayName    string `json:"display_name"`
				PickupInterval struct {
					Start time.Time `json:"start"`
					End   time.Time `json:"end"`
				} `json:"pickup_interval"`
				PickupLocation struct {
					Address struct {
						Country struct {
							IsoCode string `json:"iso_code"`
							Name    string `json:"name"`
						} `json:"country"`
						AddressLine string `json:"address_line"`
						City        string `json:"city"`
						PostalCode  string `json:"postal_code"`
					} `json:"address"`
					Location struct {
						Longitude float64 `json:"longitude"`
						Latitude  float64 `json:"latitude"`
					} `json:"location"`
				} `json:"pickup_location"`
				PurchaseEnd    time.Time `json:"purchase_end"`
				ItemsAvailable int       `json:"items_available"`
				SoldOutAt      time.Time `json:"sold_out_at"`
				Distance       float64   `json:"distance"`
				Favorite       bool      `json:"favorite"`
				InSalesWindow  bool      `json:"in_sales_window"`
				NewItem        bool      `json:"new_item"`
			} `json:"items"`
		} `json:"discover_bucket,omitempty"`
	} `json:"groupings"`
	EnabledDiscoverExperiments []string `json:"enabled_discover_experiments"`
}

func getStoresWithFood(bearerToken string) []string {

	if bearerToken == "" {
		bearerToken = getAuthBearer()
		fmt.Println("Current bearer empty, getting a new one")
	} else {
		fmt.Println("Bearer already present from last request")
	}

	requestBody := buildRequestBody()
	req, err := http.NewRequest("POST", "https://apptoogoodtogo.com/api/item/v7/discover", bytes.NewBuffer(requestBody))
	if err != nil {
		panic(err)
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", "PostmanRuntime/7.26.8")
	req.Header.Add("Accept", "*/*")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("An Error Occured %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == 401 && currentTries < MAX_TRIES {
		currentTries++
		fmt.Println("Unauthorized request, get new bearer")
		authBearer := getAuthBearer()
		getStoresWithFood(authBearer)
	} else if resp.StatusCode != 200 {
		fmt.Println("Response status: ", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var responseStruct FoodJson
	json.Unmarshal(body, &responseStruct)

	return checkStoresInResponse(responseStruct)
}

func buildRequestBody() []byte {
	userId := os.Getenv("APP_USER_ID")
	latitude := os.Getenv("LATITUDE")
	longitude := os.Getenv("LONGITUDE")

	return []byte(`{
		"user_id": "` + userId + `",
		"bucket_identifiers": ["Favorites"],
		"origin": {
			"latitude":` + latitude + `,
			"longitude":` + longitude + `
		},
		"radius": 5.0,
		"discover_experiments": ["WEIGHTED_ITEMS"]
	}`)
}

func checkStoresInResponse(response FoodJson) []string {
	var stores []string

	storesInFile := readStoresFromFile()
	for _, grouping := range response.Groupings {
		for _, item := range grouping.DiscoverBucket.Items {
			if item.ItemsAvailable > 0 {
				if !strings.Contains(storesInFile, item.Store.StoreName) {
					stores = append(stores, item.Store.StoreName)
				}
			}
		}
	}

	return stores
}
