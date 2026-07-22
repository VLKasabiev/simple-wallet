package validator

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	"github.com/shopspring/decimal"
)

var (
	letterRegex = regexp.MustCompile(`[a-zA-Z]`)
	numberRegex = regexp.MustCompile(`[0-9]`)
)

type CustomValidator struct {
	Validator *validator.Validate
}


func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func New() *CustomValidator {
	v := validator.New()

	v.RegisterValidation("strong_password", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		hasLetter := letterRegex.MatchString(password)
		hasNumber := numberRegex.MatchString(password)
		return hasLetter && hasNumber
	})
	v.RegisterValidation("decimal_gt_zero", func(fl validator.FieldLevel) bool {
        if val, ok := fl.Field().Interface().(decimal.Decimal); ok {
            return val.IsPositive()
        }
        return false
    })

	return &CustomValidator{Validator: v}
}


func FormatError(err error) map[string]string {
	errorsMap := make(map[string]string)

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		for _, f := range validationErrs {
			switch f.Tag() {
			case "required":
				errorsMap[f.Field()] = "The fiel is required"
			case "email":
				errorsMap[f.Field()] = "Invalid email address format"
			case "min":
				errorsMap[f.Field()] = fmt.Sprintf("Manimum length is %s characters", f.Param())
			case "max":
				errorsMap[f.Field()] = fmt.Sprintf("Maximum length is %s characters", f.Param())
			case "strong_password":
				errorsMap[f.Field()] = "Password must contain at least one letter and one number"
			case "decimal_gt_zero":
				errorsMap[f.Field()] = "Amount must be greater than zero"
			case "gt":
				errorsMap[f.Field()] = fmt.Sprintf("Value must be strictly greater than %s", f.Param())
			case "oneof":
				errorsMap[f.Field()] = fmt.Sprintf("Allowed values are: %s", f.Param())	
			default:
				errorsMap[f.Field()] = "Invalid value"
			}
		}
	} else {
		errorsMap["error"] = "Invalid input value"
	}

	return errorsMap
}