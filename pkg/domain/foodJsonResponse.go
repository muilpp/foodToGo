package domain

import "time"

type FoodJson struct {
	Items []struct {
		Item struct {
			ItemID     string `json:"item_id"`
			SalesTaxes []struct {
				TaxDescription string  `json:"tax_description"`
				TaxPercentage  float64 `json:"tax_percentage"`
			} `json:"sales_taxes"`
			TaxAmount struct {
				Code       string `json:"code"`
				MinorUnits int    `json:"minor_units"`
				Decimals   int    `json:"decimals"`
			} `json:"tax_amount"`
			PriceExcludingTaxes struct {
				Code       string `json:"code"`
				MinorUnits int    `json:"minor_units"`
				Decimals   int    `json:"decimals"`
			} `json:"price_excluding_taxes"`
			PriceIncludingTaxes struct {
				Code       string `json:"code"`
				MinorUnits int    `json:"minor_units"`
				Decimals   int    `json:"decimals"`
			} `json:"price_including_taxes"`
			ValueExcludingTaxes struct {
				Code       string `json:"code"`
				MinorUnits int    `json:"minor_units"`
				Decimals   int    `json:"decimals"`
			} `json:"value_excluding_taxes"`
			ValueIncludingTaxes struct {
				Code       string `json:"code"`
				MinorUnits int    `json:"minor_units"`
				Decimals   int    `json:"decimals"`
			} `json:"value_including_taxes"`
			TaxationPolicy string `json:"taxation_policy"`
			ShowSalesTaxes bool   `json:"show_sales_taxes"`
			CoverPicture   struct {
				PictureID              string `json:"picture_id"`
				CurrentURL             string `json:"current_url"`
				IsAutomaticallyCreated bool   `json:"is_automatically_created"`
			} `json:"cover_picture"`
			LogoPicture struct {
				PictureID              string `json:"picture_id"`
				CurrentURL             string `json:"current_url"`
				IsAutomaticallyCreated bool   `json:"is_automatically_created"`
			} `json:"logo_picture"`
			Name                     string        `json:"name"`
			Description              string        `json:"description"`
			FoodHandlingInstructions string        `json:"food_handling_instructions"`
			CanUserSupplyPackaging   bool          `json:"can_user_supply_packaging"`
			PackagingOption          string        `json:"packaging_option"`
			CollectionInfo           string        `json:"collection_info"`
			DietCategories           []interface{} `json:"diet_categories"`
			ItemCategory             string        `json:"item_category"`
			Buffet                   bool          `json:"buffet"`
			Badges                   []interface{} `json:"badges"`
			PositiveRatingReasons    []string      `json:"positive_rating_reasons"`
			AverageOverallRating     struct {
				AverageOverallRating float64 `json:"average_overall_rating"`
				RatingCount          int     `json:"rating_count"`
				MonthCount           int     `json:"month_count"`
			} `json:"average_overall_rating"`
			FavoriteCount int `json:"favorite_count"`
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
				PictureID              string `json:"picture_id"`
				CurrentURL             string `json:"current_url"`
				IsAutomaticallyCreated bool   `json:"is_automatically_created"`
			} `json:"logo_picture"`
			StoreTimeZone string  `json:"store_time_zone"`
			Hidden        bool    `json:"hidden"`
			FavoriteCount int     `json:"favorite_count"`
			WeCare        bool    `json:"we_care"`
			Distance      float64 `json:"distance"`
			CoverPicture  struct {
				PictureID              string `json:"picture_id"`
				CurrentURL             string `json:"current_url"`
				IsAutomaticallyCreated bool   `json:"is_automatically_created"`
			} `json:"cover_picture"`
			IsManufacturer bool `json:"is_manufacturer"`
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
		Distance       float64   `json:"distance"`
		Favorite       bool      `json:"favorite"`
		InSalesWindow  bool      `json:"in_sales_window"`
		NewItem        bool      `json:"new_item"`
	} `json:"items"`
}
