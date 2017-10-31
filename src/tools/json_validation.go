package tools

import (
	"fmt"
	"gopkg.in/go-playground/validator.v9"
	"models"
	"reflect"
)

type DefaultValidator struct {
	validate *validator.Validate
}

func (v *DefaultValidator) ValidateIncomingJsonRequest(obj interface{}) models.Error {
	ValidationError := models.Error{
		"SUCCESS", 200, "Successful request", map[string][]string{}, nil, "Successful request",
	}
	if KindOfData(obj) == reflect.Struct {
		v.lazyinit()
		if err := v.validate.Struct(obj); err != nil {
			ValidationError.Id = "bad_request"
			ValidationError.Msg = "Ohy Vey!"
			ValidationError.Title = "Bad request"
			ValidationError.Status = 400
			for _, fieldErr := range err.(validator.ValidationErrors) {
				errorString := fmt.Sprintf("Failed %s validation", fieldErr.Tag())
				ValidationError.Details[fieldErr.StructField()] = append(
					ValidationError.Details[fieldErr.StructField()],
					errorString,
				)
			}
		}
	}
	return ValidationError
}

func (v *DefaultValidator) lazyinit() {
	v.validate = validator.New()
}

func KindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}
