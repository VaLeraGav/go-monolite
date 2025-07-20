package validator

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	customFieldPhone        = "phone"
	customFieldAuthPassword = "auth_password"
)

func initCustomValidPatterns() error {
	if err := validatorIns.RegisterValidation(customFieldPhone, createValidPhone); err != nil {
		return err
	}

	if err := validatorIns.RegisterValidation(customFieldAuthPassword, createValidAuthPasswordPattern); err != nil {
		return err
	}

	return nil
}

func createValidPattern(fl validator.FieldLevel, pattern string) bool {
	return regexp.MustCompile(pattern).MatchString(fl.Field().String())
}

func createValidPhone(fl validator.FieldLevel) bool {
	return createValidPattern(fl, `^7\d{10}$`)
}

func createValidAuthPasswordPattern(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString(value) // Есть хотя бы одна заглавная буква
	hasLowercase := regexp.MustCompile(`[a-z]`).MatchString(value) // Есть хотя бы одна строчная буква
	hasDigit := regexp.MustCompile(`\d`).MatchString(value)        // Есть хотя бы одна цифра

	return hasUppercase && hasLowercase && hasDigit
}
