package restful

import (
	"Infinite_train/pkg/common/constant"
	"Infinite_train/pkg/common/utils/log/golog"
	"github.com/dlclark/regexp2"
	"gopkg.in/go-playground/validator.v9"
)

func CheckPort(fl validator.FieldLevel) bool {
	if fl.Field().Int() > 0 || fl.Field().Int() < 65536 {
		return true
	}
	return false
}

func InstanceName(fl validator.FieldLevel) bool {
	if fl.Field().Len() > 16 || fl.Field().Len() < 2 {
		return false
	}
	re, err := regexp2.Compile(constant.InstanceNameRegEx, 0)
	if err != nil {
		return false
	}
	result, err := re.FindStringMatch(fl.Field().String())
	if err != nil || result == nil {
		return false
	}
	return true
}

func NewCustomValidator() *validator.Validate {
	validate := validator.New()

	var err error
	err = validate.RegisterValidation("CheckPort", CheckPort)
	if err != nil {
		golog.Warnf("0", "validate Register CheckPort failed: %s", err.Error())
	}
	err = validate.RegisterValidation("InstanceName", InstanceName)
	if err != nil {
		golog.Warnf("0", "validate Register InstanceName failed: %s", err.Error())
	}
	return validate
}
