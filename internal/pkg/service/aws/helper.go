package aws

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer"
	"github.com/aws/aws-sdk-go-v2/service/costexplorer/types"
	"strconv"
)

var (
	groupByDimension = func(dimensions []string) []types.GroupDefinition {
		var groups []types.GroupDefinition
		for _, d := range dimensions {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(d),
			})
		}
		return groups
	}
	groupByTag = func(tag string) []types.GroupDefinition {
		return []types.GroupDefinition{
			{
				Type: types.GroupDefinitionTypeTag,
				Key:  aws.String(tag),
			},
		}
	}
	groupByTagAndDimension = func(tag string, dimensions []string) []types.GroupDefinition {
		var groups []types.GroupDefinition
		for _, d := range dimensions {
			groups = append(groups, types.GroupDefinition{
				Type: types.GroupDefinitionTypeDimension,
				Key:  aws.String(d),
			})
		}
		groups = append(groups, types.GroupDefinition{
			Type: types.GroupDefinitionTypeTag,
			Key:  aws.String(tag),
		})
		return groups
	}
	filterCredits = func() *types.Expression {
		return &types.Expression{
			Not: &types.Expression{
				Dimensions: &types.DimensionValues{
					Key:    "RECORD_TYPE",
					Values: []string{"Refund", "Credit", "DiscountedUsage"},
				},
			},
		}
	}
	filterByTag = func(tag string, value string) *types.Expression {
		return &types.Expression{
			Tags: &types.TagValues{
				Key:    aws.String(tag),
				Values: []string{value},
			},
		}
	}
)

func ConvertToFloat(amount string) float64 {
	f, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func ReturnIfPresent(s []string) string {
	if len(s) == 1 {
		return ""
	} else {
		return s[1]
	}

}

func ToSlice(d costexplorer.GetDimensionValuesOutput) []string {
	var servicesSlice []string
	for _, service := range d.DimensionValues {
		servicesSlice = append(servicesSlice, *service.Value)
	}
	return servicesSlice
}

func CostAndUsageFilterGenerator(req CostAndUsageRequestType) *types.
	Expression {
	expression := &types.Expression{}

	if req.ExcludeDiscounts && req.IsFilterEnabled {
		expression.And = []types.Expression{*filterCredits(),
			*filterByTag(req.Tag, req.TagFilterValue)}
	} else if req.ExcludeDiscounts {
		expression = filterCredits()
	} else if req.IsFilterEnabled {
		expression = filterByTag(req.Tag, req.TagFilterValue)
	} else {
		return nil
	}
	return expression
}

func CostForecastFilterGenerator(req GetCostForecastRequest) *types.
	Expression {
	var filterExpression types.Expression
	var expList []types.Expression
	var exp types.Expression

	var isMultiFilter bool
	if len(req.Filter.Dimensions) > 1 {
		isMultiFilter = true
	}

	for _, dimension := range req.Filter.Dimensions {
		temp := &types.DimensionValues{
			Key:    types.Dimension(dimension.Key),
			Values: dimension.Value,
		}

		if len(req.Filter.Dimensions) == 1 {
			expList = append(expList, types.Expression{
				Dimensions: temp,
			})
		} else if len(req.Filter.Dimensions) > 1 {
			exp.And = append(exp.And, types.Expression{
				Dimensions: temp,
			})
		}
	}

	if isMultiFilter {
		expList = append(expList, exp)
	}

	//
	//for _, tag := range req.Filter.Tags {
	//	temp := &types.TagValues{
	//		Key:    aws.String(tag.Key),
	//		Values: tag.Value,
	//	}
	//	exp = append(exp, types.Expression{
	//		Tags: temp,
	//	})
	//}

	filterExpression = expList[0]

	return &filterExpression
}

func CostAndUsageGroupByGenerator(req CostAndUsageRequestType) []types.GroupDefinition {
	if req.Tag != "" && len(req.GroupBy) == 1 {
		return groupByTagAndDimension(req.Tag, req.GroupBy)
	} else if req.Tag != "" {
		return groupByTag(req.Tag)
	} else {
		return groupByDimension(req.GroupBy)
	}
}
