package validator

import "errors"

var (
	ErrorValidation = errors.New("ошибка в валидации поля")
	ErrorRequire    = errors.New("обязательно для заполнения")
)
