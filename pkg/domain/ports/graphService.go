package ports

const FOOD_CHART_BY_STORE = "food-by-store.gif"
const FOOD_CHART_BY_HOUR_OF_DAY = "food-by-hour.gif"
const FOOD_CHART_BY_DAY_OF_WEEK = "food-by-day.gif"

type GraphService interface {
	PrintAllReports()
}
