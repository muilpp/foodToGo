package ports

const FOOD_CHART_BY_STORE = "food-chart-store.png"
const FOOD_CHART_BY_HOUR_OF_DAY = "food-chart-hour.png"
const FOOD_CHART_BY_DAY_OF_WEEK = "food-chart-week.png"

type GraphService interface {
	PrintAllReports()
}
