package writers

import (
	"github.com/cduggn/ccexplorer/internal/core/domain/model"
	"os"
	"strings"
)

var (
	maxDisplayRows = 10
	csvFileName    = "ccexplorer.csv"
	csvHeader      = []string{"Dimension/Tag", "Dimension/Tag", "Metric",
		"Granularity",
		"Start",
		"End", "USD Amount", "Unit"}
	OutputDir = "./output"
)

type Builder struct {
}

func CostAndUsageToStdout(sortFn func(r map[int]model.Service) []model.Service,
	r model.CostAndUsageOutputType) error {

	sortedServices := sortFn(r.Services)
	output := ConvertToStdoutType(sortedServices, r.Granularity)

	w, err := NewStdoutWriter("costAndUsage")
	if err != nil {
		return model.Error{
			Msg: "Error writing to stdout : " + err.Error()}
	}
	w.Writer(output)
	return nil
}

func CostAndUsageToCSV(sortFn func(r map[int]model.Service) []model.Service,
	r model.CostAndUsageOutputType) error {

	f, err := NewCSVFile(OutputDir, csvFileName)
	if err != nil {
		return model.Error{
			Msg: "Error creating CSV file: " + err.Error()}
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	rows := ConvertServiceMapToArray(r.Services, r.Granularity)
	err = WriteToCSV(f, csvHeader, rows)
	if err != nil {
		return model.Error{
			Msg: "Error writing to CSV file: " + err.Error()}
	}
	return nil
}

func CostAndUsageToChart(sortFn func(r map[int]model.Service) []model.Service,
	r model.CostAndUsageOutputType) error {

	builder := Builder{}
	s := sortFn(r.Services)
	input := ConvertToChartInputType(r, s)

	charts, err := builder.NewCharts(input)
	if err != nil {
		return err
	}

	err = WriteToChart(charts)
	if err != nil {
		return err
	}
	return nil
}

func CostAndUsageToOpenAI(sortFn func(r map[int]model.Service) []model.Service,
	r model.CostAndUsageOutputType) error {

	sortedData := sortFn(r.Services)
	rows := ConvertServiceSliceToArray(sortedData, r.Granularity)

	maxRows := MaxSupportedRows(rows, maxDisplayRows)
	data := BuildPromptText(rows[:maxRows])
	resp, err := Summarize(r.OpenAIAPIKey, data)
	if err != nil {
		return err
	}
	err = WriteToHTML(resp.Choices[0].Message.Content)

	if err != nil {
		return err
	}
	return nil
}

func MaxSupportedRows(rows [][]string, maxRows int) int {
	if len(rows) > maxRows {
		return maxRows
	}
	return len(rows)
}

func ForecastToStdout(r model.ForecastPrintData,
	dimensions []string) {

	filteredBy := strings.Join(dimensions, " | ")
	output := ConvertToForecastStdoutType(r, filteredBy)
	w, err := NewStdoutWriter("forecast")
	if err != nil {
		return
	}
	w.Writer(output)
}
