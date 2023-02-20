package printer

import (
	"bytes"
	"context"
	"fmt"
	gogpt "github.com/sashabaranov/go-gpt3"
	"html/template"
	"os"
	"strings"
)

type TrainingData struct {
	Dimension   string
	Tag         string
	Metric      string
	Granularity string
	Start       string
	End         string
	USDAmount   string
	Unit        string
}

var (
	aiFileName       = "ccexplorer_ai.html"
	trainingTemplate = `
        <table>
            <thead>
                <tr>
                    <th>Dimension/Tag</th>
                    <th>Dimension/Tag</th>
                    <th>Start</th>
                    <th>End</th>
                    <th>USD Amount</th>
                </tr>
            </thead>
            <tbody>
                {{range .}}
                <tr>
                    <td>{{.Dimension}}</td>
                    <td>{{.Tag}}</td>
                    <td>{{.Start}}</td>
                    <td>{{.End}}</td>
                    <td>{{.USDAmount}}</td>
                </tr>
                {{end}}
            </tbody>
        </table>
    `
)

func AIWriter(f *os.File, completions string) error {
	_, err := f.WriteString(completions)
	if err != nil {
		return PrinterError{
			msg: "Failed writing to AI HTML: " + err.Error(),
		}
	}
	return nil
}

func ConvertToCommaDelimitedString(rows [][]string) string {
	var buf bytes.Buffer

	buf.WriteString(csvHeaderPromptFormat)

	for i, row := range rows {
		for j, col := range row {
			buf.WriteString(col)
			if j < len(row)-1 {
				buf.WriteByte(',')
			}
		}
		if i < len(rows)-1 {
			buf.WriteByte(';')
		}
	}
	cvsString := buf.String()
	return cvsString
}

func BuildPromptText(rows [][]string) string {
	var builder strings.Builder
	builder.WriteString("Generate a html table that looks like this ")

	trainingData := BuildCostAndUsagePromptText(rows)
	builder.WriteString(trainingData)

	builder.WriteString(" from the following csv data, " +
		"showing top 10 rows ")
	costAndUsageData := ConvertToCommaDelimitedString(rows[:20])
	builder.WriteString(costAndUsageData)

	builder.WriteString(" Display the title Cost And Usage Report above" +
		" the table in h2 font ," +
		"display USD currency and date range in h3 font below the title. " +
		"Style the table with css")

	builder.WriteString("Place a hr tag below the table, below" +
		"generate recommendations in bullet list for cost optimization for" +
		" most expensive rows. Recommendations should resemble this below but" +
		" using the table data")

	builder.WriteString("Health-Check-Option-AWS => Use Amazon Route 53's Traffic Flow feature: Amazon Route 53's Traffic Flow feature allows you to route traffic based on endpoint health")

	return builder.String()
}

func BuildCostAndUsagePromptText(rows [][]string) string {
	t := CreateTrainingTemplate()
	s, err := CreateTrainingData(t, BuildTrainingDataRow(rows))
	if err != nil {
		fmt.Println("Error populating template: ", err)
	}
	return s
}

func CreateTrainingTemplate() *template.Template {
	t := template.Must(template.New("table").Parse(trainingTemplate))
	return t
}

func CreateTrainingData(t *template.Template, data []TrainingData) (string,
	error) {
	var buf bytes.Buffer
	err := t.Execute(&buf, data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func BuildTrainingDataRow(rows [][]string) []TrainingData {
	return []TrainingData{
		{
			Dimension:   rows[0][0],
			Tag:         rows[0][1],
			Metric:      rows[0][2],
			Granularity: rows[0][3],
			Start:       rows[0][4],
			End:         rows[0][5],
			USDAmount:   rows[0][6],
			Unit:        rows[0][7],
		},
	}
}

func SummarizeWIthAI(apiKey string, promptData string) (gogpt.
	CompletionResponse,
	error) {

	fmt.Println("Generating costAndUsage report with gpt3...")

	c := gogpt.NewClient(apiKey)
	ctx := context.Background()

	req := gogpt.CompletionRequest{
		Model:     gogpt.GPT3TextDavinci003,
		MaxTokens: 2400,
		Prompt:    promptData,
	}
	resp, err := c.CreateCompletion(ctx, req)
	if err != nil {
		return gogpt.CompletionResponse{}, err
	}

	return resp, nil
}
