package aws

import (
	"time"
)

type ValidationError struct {
	msg string
}

func (e ValidationError) Error() string {
	return e.msg
}

//func ValidateGroupByMap(groupBy map[string]string) ([]string, error) {
//
//	var tag map[string]string
//
//	for key, val := range groupBy {
//		if val == "DIMENSION" {
//			//dimension = map[string]string{key: val}
//			err := ValidateGroupByDimension(val)
//			if err != nil {
//				return nil, err
//			}
//		}
//		if val == "TAG" {
//			tag = map[string]string{key: val}
//			//ValidateGroupByTag(tag)
//		} else {
//			return nil, ValidationError{
//				msg: "GroupBy must be one of the following: DIMENSION, TAG",
//			}
//		}
//	}
//	return nil, nil
//}

func ValidateGroupByDimension(dimension string) error {

	switch dimension {
	case "AZ", "SERVICE", "USAGE_TYPE", "INSTANCE_TYPE", "LINKED_ACCOUNT", "OPERATION",
		"PURCHASE_TYPE", "PLATFORM", "TENANCY", "RECORD_TYPE", "LEGAL_ENTITY_NAME",
		"INVOICING_ENTITY", "DEPLOYMENT_OPTION", "DATABASE_ENGINE",
		"CACHE_ENGINE", "INSTANCE_TYPE_FAMILY", "REGION", "BILLING_ENTITY",
		"RESERVATION_ID", "SAVINGS_PLANS_TYPE", "SAVINGS_PLAN_ARN",
		"OPERATING_SYSTEM":
	default:
		return ValidationError{
			msg: "Dimension must be one of the following: AZ, SERVICE, " +
				"USAGE_TYPE, INSTANCE_TYPE, LINKED_ACCOUNT, OPERATION, " +
				"PURCHASE_TYPE, PLATFORM, TENANCY, RECORD_TYPE, " +
				"LEGAL_ENTITY_NAME, INVOICING_ENTITY, DEPLOYMENT_OPTION, " +
				"DATABASE_ENGINE, CACHE_ENGINE, INSTANCE_TYPE_FAMILY, " +
				"REGION, BILLING_ENTITY, RESERVATION_ID, " +
				"SAVINGS_PLANS_TYPE, SAVINGS_PLAN_ARN, OPERATING_SYSTEM",
		}
	}

	return nil
}

func ValidateStartDate(startDate string) error {
	if startDate == "" {
		return ValidationError{
			msg: "Start date must be specified",
		}
	}

	start, _ := time.Parse("2006-01-02", startDate)
	today := time.Now()
	if start.After(today) {
		return ValidationError{
			msg: "Start date must be before today's date",
		}
	}

	return nil
}

func ValidateEndDate(endDate, startDate string) error {
	if endDate == "" {
		return ValidationError{
			msg: "End date must be specified",
		}
	}

	end, _ := time.Parse("2006-01-02", endDate)
	today := time.Now()
	if end.After(today) {
		return ValidationError{
			msg: "End date must be before today's date",
		}
	}

	start, _ := time.Parse("2006-01-02", startDate)
	if end.Before(start) {
		return ValidationError{
			msg: "End date must not be before start date",
		}
	}

	return nil
}
