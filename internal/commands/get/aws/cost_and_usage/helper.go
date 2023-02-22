package cost_and_usage

func IsValidPrintFormat(f string) bool {
	return f == "stdout" || f == "csv" || f == "chart" || f == "gpt3"
}

func IsValidGranularity(g string) bool {
	return g == "DAILY" || g == "MONTHLY" || g == "HOURLY"
}

func IsValidMetric(m string) bool {
	return m == "AmortizedCost" || m == "BlendedCost" || m == "NetAmortizedCost" ||
		m == "NetUnblendedCost" || m == "NormalizedUsageAmount" || m == "UnblendedCost" ||
		m == "UsageQuantity"
}

func SortByFn(sortByDate bool) string {
	if sortByDate {
		return "date"
	}
	return "cost"
}

func HasAccountInformation(groupBy []string) bool {
	for _, v := range groupBy {
		if v == "LINKED_ACCOUNT" {
			return true
		}
	}
	return false

}
