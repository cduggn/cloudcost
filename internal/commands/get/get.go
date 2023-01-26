package get

import (
	"github.com/cduggn/ccexplorer/internal/commands/get/aws"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
)

var (
	costAndUsageCmd = &cobra.Command{
		Use:   "get",
		Short: "Cost and usage summary for AWS services",
		Long:  paintHeader(),
	}
	awsCost = &cobra.Command{
		Use:   "aws",
		Short: "Explore UNBLENDED cost summaries for AWS",
		Long: `
Command: aws 
Description: Cost and usage summary for AWS services.

Prerequisites:
- AWS credentials configured in ~/.aws/credentials and default region configured in ~/.aws/config. Alternatively, 
you can set the environment variables AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY and AWS_REGION.`,
		Example: costAndUsageExamples,
		RunE:    aws.CostAndUsageSummary,
	}
	forecast = &cobra.Command{
		Use:     "forecast",
		Short:   "Return cost and usage forecasts for your account.",
		Example: forecastExamples,
		RunE:    aws.CostForecast,
	}
)

func paintHeader() string {
	myFigure := figure.NewFigure("Cost And Usage", "thin", true)
	return myFigure.String()
}

func AWSCostAndUsageCommand() *cobra.Command {

	costAndUsageCmd.AddCommand(aws.CostAndUsageCommand(awsCost))
	awsCost.AddCommand(aws.ForecastCommand(forecast))
	return costAndUsageCmd
}
