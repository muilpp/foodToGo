package ports

const FOOD_CHART_BY_STORE_MONTHLY = "-food-by-store-monthly.gif"
const FOOD_CHART_BY_HOUR_OF_DAY_MONTHLY = "-food-by-hour-monthly.gif"
const FOOD_CHART_BY_DAY_OF_WEEK_MONTHLY = "-food-by-day-monthly.gif"

const FOOD_CHART_BY_STORE_YEARLY = "-food-by-store-yearly.gif"
const FOOD_CHART_BY_HOUR_OF_DAY_YEARLY = "-food-by-hour-yearly.gif"
const FOOD_CHART_BY_DAY_OF_WEEK_YEARLY = "-food-by-day-yearly.gif"

type GraphService interface {
	PrintAllMonthlyReports()
	PrintAllYearlyReports()
}
