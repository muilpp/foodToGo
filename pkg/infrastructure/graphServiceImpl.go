package infrastructure

import (
	"os"
	"sort"
	"strconv"
	"time"

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

func (gs GraphServiceImpl) PrintAllReports() {
	result := gs.repository.GetStoresByHourOfDay()
	gs.printValuesToFile(result, ports.FOOD_CHART_BY_HOUR_OF_DAY, "(by hour of day)")

	result = gs.repository.GetStoresByDayOfWeek()
	gs.printValuesToFile(result, ports.FOOD_CHART_BY_DAY_OF_WEEK, "(by day of week)")

	result = gs.repository.GetStoresByTimesAppeared()
	gs.printValuesToFile(result, ports.FOOD_CHART_BY_STORE, "(by store)")
}

func (gs GraphServiceImpl) printValuesToFile(storeCounter []domain.StoreCounter, fileName string, subtitle string) {
	valueSlice := getValuesToPlot(storeCounter)
	tickSlice := getYAxisLabels(float64(valueSlice[0].Value))

	year, month, _ := time.Now().AddDate(0, 0, -1).Date()
	title := month.String() + " " + strconv.Itoa(year) + " " + subtitle

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

func getYAxisLabels(max float64) []chart.Tick {
	var tickSlice []chart.Tick

	for i := 0; i <= int(max); i++ {
		tickSlice = append(tickSlice, chart.Tick{Value: float64(i), Label: strconv.Itoa(int(i))})
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
