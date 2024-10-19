package helpers

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func GetValidatorError(validatorError error, structRequest interface{}) (errorMessages []ErrorMessage) {
	validationErrors := validatorError.(validator.ValidationErrors)
	val := reflect.ValueOf(structRequest)
	for _, fieldError := range validationErrors {
		var errorMessage ErrorMessage
		structField, ok := val.Type().FieldByName(fieldError.Field())
		if !ok {
			errorMessage.Field = "property"
			errorMessage.Message = "couldn't find property: " + fieldError.Field()
			errorMessages = append(errorMessages, errorMessage)
			return
		}
		errorMessage.Field = structField.Tag.Get("json")
		if fieldError.Tag() == "usernamevalidator" {
			errorMessage.Message = "please use only uppercase and lowercase letter and number and min 5 and max 8 alphanumeric"
		} else if fieldError.Tag() == "passwordvalidator" {
			errorMessage.Message = "please use only uppercase, lowercase, number and must have 1 uppercase. lowercase, number, @, _, -, min 8 and max 20"
		} else if fieldError.Tag() == "telephonevalidator" {
			errorMessage.Message = "please use only number and + "
		} else if fieldError.Tag() == "email" {
			errorMessage.Message = "please input a correct email format "
		} else if fieldError.Tag() == "gte" {
			errorMessage.Message = "please input greater than equal to " + fieldError.Param()
		} else {
			errorMessage.Message = "is " + fieldError.Tag()
		}
		errorMessages = append(errorMessages, errorMessage)
	}
	return
}
