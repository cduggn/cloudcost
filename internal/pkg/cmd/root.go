package cmd

import (
	"github.com/cduggn/cloudcost/internal/pkg/cmd/billing"
	"github.com/common-nighthawk/go-figure"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "cloudcost",
	Short: "A CLI tool to get AWS Costs, Usage and Forecasts",
	Long:  paintHeader(),
}

func paintHeader() string {
	myFigure := figure.NewFigure("CloudCost", "thin", true)
	return myFigure.String()
}

func init() {
	rootCmd.AddCommand(billing.CostAndUsageCommand())
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(126)
	}
}
