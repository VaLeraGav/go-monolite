package validator

import (
	"errors"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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

func initCustomValidPatterns() error {
	if err := validatorIns.RegisterValidation("phone", createValidPhone); err != nil {
		return err
	}

	if err := validatorIns.RegisterValidation("auth_password", createValidPasswordPattern); err != nil {
		return err
	}

	if err := validatorIns.RegisterValidation("auth_password", createValidPasswordPattern); err != nil {
		return err
	}

	return nil
}

func CompareFields(f1 any, f2 any, tag string) error {
	return validatorIns.VarWithValue(f1, f2, tag)
}

func Validate(dto any, tagMessages map[string]string) error {
	err := validatorIns.Struct(dto)

	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	ve := getErrorsMap(validationErrors, tagMessages)

	if len(ve.Fields) == 0 {
		return nil
	}
	return ve
}

func ParseUUID(uuidStr string) (uuid.UUID, error) {
	if uuidStr == "" {
		return uuid.UUID{}, errors.New("UUID не указан")
	}

	uuidParse, err := uuid.Parse(uuidStr)
	if err != nil {
		return uuid.UUID{}, errors.New("некорректный формат UUID")
	}

	return uuidParse, nil
}

func getErrorsMap(validationErrors validator.ValidationErrors, tagMessages map[string]string) ValidationError {
	errs := make(map[string]string)

	for _, fieldError := range validationErrors {
		field := strings.ToLower(fieldError.Field())

		if msg, exists := tagMessages[fieldError.Tag()]; exists {
			errs[field] = "Поле " + field + " " + msg
		} else {
			errs[field] = "Поле " + field + " содержит ошибку валидации: " + fieldError.Tag()
		}
	}

	if len(errs) == 0 {
		return ValidationError{}
	}

	return ValidationError{Err: validationErrors, Fields: errs}
}

func createValidPattern(fl validator.FieldLevel, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(fl.Field().String())
}

func createValidPasswordPattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(value)
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(value)
	hasDigit := regexp.MustCompile(`\d`).MatchString(value)

	return hasUppercase && hasLowercase && hasDigit
}

func createValidPhone(fl validator.FieldLevel) bool {
	return createValidPattern(fl, `^7\d{10}$`)
}
