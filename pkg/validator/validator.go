package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Err    error
	Fields map[string]string
}

func (v ValidationError) Error() string {
	return v.Err.Error()
}

var validatorIns *validator.Validate

func init() {
	validatorIns = validator.New()
	if err := initCustomValidPatterns(); err != nil {
		panic(err)
	}
}

func CompareFields(f1 any, f2 any, tag string) error {
	return validatorIns.VarWithValue(f1, f2, tag)
}

func Validate(dto any) error {
	return ValidateWithMessages(dto, map[string]string{})
}

func ValidateWithMessages(dto any, tagMessages map[string]string) error {
	err := validatorIns.Struct(dto)
	if err == nil {
		return nil
	}

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}

	ve := getErrorsMap(validationErrors, tagMessages)
	if len(ve.Fields) == 0 {
		return nil
	}

	return ve
}

func getErrorsMap(validationErrors validator.ValidationErrors, tagMessages map[string]string) ValidationError {
	errors := make(map[string]string)

	for _, err := range validationErrors {
		tag := err.Tag()
		field := ToSnakeCase(err.Field())

		msg, ok := tagMessages[tag]
		if !ok {
			msg = defaultErrorMessage(tag, field, err.Param())
		}

		errors[field] = msg
	}

	if len(errors) == 0 {
		return ValidationError{}
	}

	return ValidationError{Err: validationErrors, Fields: errors}
}

func defaultErrorMessage(tag, field, param string) string {
	switch tag {
	case "required":
		return fmt.Sprintf("Поле %s обязательно для заполнения", field)
	case "min":
		return fmt.Sprintf("Поле %s должно быть не меньше %s", field, param)
	case "max":
		return fmt.Sprintf("Поле %s должно быть не больше %s", field, param)
	case "email":
		return fmt.Sprintf("Поле %s должно быть валидным email адресом", field)
	case "gte":
		return fmt.Sprintf("Поле %s должно быть больше или равно %s", field, param)
	case "lte":
		return fmt.Sprintf("Поле %s должно быть меньше или равно %s", field, param)
	case "eqfield":
		return fmt.Sprintf("Поле %s должно совпадать с другим полем", field)
	case "oneof":
		return fmt.Sprintf("Поле %s должно быть одним из следующих значений: %s", field, param)
	case customFieldPhone:
		return fmt.Sprintf("Поле %s неправильно передано номер телефона", field)
	case customFieldAuthPassword:
		return fmt.Sprintf("Поле %s неправильно пароль", field)
	default:
		return fmt.Sprintf("Поле %s не прошло проверку: %s", field, tag)
	}
}

func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
