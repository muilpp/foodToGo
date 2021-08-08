package infrastructure

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/jinzhu/now"
	"github.com/marc/get-food-to-go/pkg/domain"
	"github.com/marc/get-food-to-go/pkg/domain/ports"
	"github.com/wcharczuk/go-chart/v2"
)

type GraphServiceImpl struct {
	repository ports.Repository
}

func NewGraphService(repository ports.Repository) *GraphServiceImpl {
	return &GraphServiceImpl{repository}
}

func (gs GraphServiceImpl) PrintAllMonthlyReports() {
	lastMonth := now.BeginningOfMonth().AddDate(0, -1, 0).Format("2006-01-02")

	result := gs.repository.GetStoresByHourOfDay(lastMonth)
	gs.printValuesToFile(result, ports.FOOD_CHART_BY_HOUR_OF_DAY_MONTHLY, "(by hour of day)", false)

	result = gs.repository.GetStoresByDayOfWeek(lastMonth)
	gs.printValuesToFile(result, ports.FOOD_CHART_BY_DAY_OF_WEEK_MONTHLY, "(by day of week)", false)

	result = gs.repository.GetStoresByTimesAppeared(lastMonth)
	gs.printValuesToFile(result, ports.FOOD_CHART_BY_STORE_MONTHLY, "(by store)", false)
}

func (gs GraphServiceImpl) PrintAllYearlyReports() {
	lastYear := now.BeginningOfYear().AddDate(0, -1, 0).Format("2006-01-02")
	fmt.Println("Last year: ", lastYear)
	result := gs.repository.GetStoresByHourOfDay(lastYear)

	fmt.Println("Values found: ", len(result))

	gs.printValuesToFile(result, ports.FOOD_CHART_BY_HOUR_OF_DAY_YEARLY, "(by hour of day)", true)

	result = gs.repository.GetStoresByDayOfWeek(lastYear)
	gs.printValuesToFile(result, ports.FOOD_CHART_BY_DAY_OF_WEEK_YEARLY, "(by day of week)", true)

	result = gs.repository.GetStoresByTimesAppeared(lastYear)
	gs.printValuesToFile(result, ports.FOOD_CHART_BY_STORE_YEARLY, "(by store)", true)
}

func (gs GraphServiceImpl) printValuesToFile(storeCounter []domain.StoreCounter, fileName string, subtitle string, isYearlyReport bool) {
	valueSlice := getValuesToPlot(storeCounter)
	tickSlice := getYAxisLabels(storeCounter)

	year, month, _ := time.Now().AddDate(0, 0, -1).Date()

	var title string
	if isYearlyReport {
		title = strconv.Itoa(year) + " report " + subtitle
	} else {
		title = month.String() + " " + strconv.Itoa(year) + " " + subtitle
	}

	graph := chart.BarChart{
		Title: title,
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars:     valueSlice,

		YAxis: chart.YAxis{
			Ticks: tickSlice,
		},
	}
	f, _ := os.Create(fileName)
	defer f.Close()
	graph.Render(chart.PNG, f)
}

func getYAxisLabels(storeSlice []domain.StoreCounter) []chart.Tick {
	var tickSlice []chart.Tick

	tickSlice = append(tickSlice, chart.Tick{Value: float64(0), Label: strconv.Itoa(int(0))})
	for _, v := range storeSlice {
		tickSlice = append(tickSlice, chart.Tick{Value: float64(v.GetTotal()), Label: strconv.Itoa(int(v.GetTotal()))})
	}

	return tickSlice
}

func getValuesToPlot(result []domain.StoreCounter) chart.Values {

	var valueSlice chart.Values
	for _, v := range result {
		valueSlice = append(valueSlice, chart.Value{Value: float64(v.GetTotal()), Label: v.GetName()})
	}

	sort.Slice(valueSlice, func(i, j int) bool {
		return valueSlice[i].Value > valueSlice[j].Value
	})

	return valueSlice
}
