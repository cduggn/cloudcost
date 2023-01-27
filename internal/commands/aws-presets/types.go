package aws_presets

type PresetError struct {
	msg string
}

type PresetParams struct {
	Alias             string
	Dimension         []string
	Tag               string
	Filter            map[string]string
	FilterType        string
	FilterByDimension bool
	FilterByTag       bool
	ExcludeDiscounts  bool
	CommandSyntax     string
	Description       []string
}
