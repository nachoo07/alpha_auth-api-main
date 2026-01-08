// Package validator contains a helper to centralize the golang struct validator global reference
package validator

import (
	"regexp"
	"time"

	"github.com/go-playground/validator/v10"
)

type ValidationErrors = validator.ValidationErrors

var v = validator.New()

func Validate(entity any) error {
	return v.Struct(entity)
}

func IsValidTimeFormat(time string) bool {
	if time == "" {
		return true
	}

	match, err := regexp.MatchString(`^(0[0-9]|1[0-9]|2[0-3]):([0-5][0-9])$`, time)

	return err == nil && match
}

func IsValidChannelFormat(channelType string) bool {
	match, err := regexp.MatchString(`^(logistics|commercial)$`, channelType)

	return err == nil && match
}

func IsValidStepTypeFormat(stepType string) bool {
	match, err := regexp.MatchString(`^(first|last|middle)_mile$`, stepType)

	return err == nil && match
}

func IsValidModalFormat(modal string) bool {
	match, err := regexp.MatchString(`^(air|ground)$`, modal)

	return err == nil && match
}

func IsValidSiteFormat(siteID string) bool {
	match, err := regexp.MatchString(`^([A-Z]){3}$`, siteID)

	return err == nil && match
}

// IsValidDateFormat is a function thats checks if a string is valid following these rules:
// - YYYY-MM-DD
// - YYYY min 0000 and max 9999
// - MM min 01 and max 12
// - DD min 01 and max 31
// - Example: 2020-12-31.
// Note: Empty values "" will be considered valid, since the required tag must be used to check if mandatory.
func IsValidDateFormat(date string) bool {
	if date == "" {
		return true
	}

	_, err := time.Parse(time.DateOnly, date)
	return err == nil
}

// IsDateFromBeforeOrEqualDateTo is a function that checks if the dates, in string UTC format, where dateFrom occurs before or equal dateTo
func IsDateFromBeforeOrEqualDateTo(dateFrom, dateTo time.Time) bool {
	return dateFrom.Before(dateTo) || dateFrom.Equal(dateTo)
}

var validationMap = map[string]func(validator.FieldLevel) bool{
	"validSiteFormat":     validateSiteFormat,
	"validTimeFormat":     validateTimeFormat,
	"validStepTypeFormat": validateStepTypeFormat,
	"validChannelFormat":  validateChannelFormat,
	"validModalFormat":    validateModalFormat,
	"validDateFormat":     validateDateFormat,
}

var registered bool

func RegisterValidation() (err error) {
	if !registered {
		err = registerValidationMap(validationMap)
		registered = true
	}
	return err
}

func registerValidationMap(m map[string]func(validator.FieldLevel) bool) error {
	for key, value := range m {
		if err := v.RegisterValidation(key, value); err != nil {
			return err
		}
	}

	return nil
}

func validateSiteFormat(fl validator.FieldLevel) bool {
	return IsValidSiteFormat(fl.Field().String())
}

func validateTimeFormat(fl validator.FieldLevel) bool {
	return IsValidTimeFormat(fl.Field().String())
}

func validateChannelFormat(fl validator.FieldLevel) bool {
	return IsValidChannelFormat(fl.Field().String())
}

func validateStepTypeFormat(fl validator.FieldLevel) bool {
	return IsValidStepTypeFormat(fl.Field().String())
}

func validateModalFormat(fl validator.FieldLevel) bool {
	return IsValidModalFormat(fl.Field().String())
}

func validateDateFormat(fl validator.FieldLevel) bool {
	return IsValidDateFormat(fl.Field().String())
}
