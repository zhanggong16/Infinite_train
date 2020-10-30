package common

import (
	"Infinite_train/pkg/common/utils/log/golog"
	"gopkg.in/go-playground/validator.v9"
)

func CheckPort(fl validator.FieldLevel) bool {
	if fl.Field().Int() > 0 || fl.Field().Int() < 65536 {
		return true
	}
	return false
}

func NewCustomValidator() *validator.Validate {
	validate := validator.New()

	var err error
	err = validate.RegisterValidation("CheckPort", CheckPort)
	if err != nil {
		golog.Warnf("0", "validate Register CheckPort failed: %s", err.Error())
	}

	return validate
}