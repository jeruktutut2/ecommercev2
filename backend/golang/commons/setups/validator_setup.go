package setups

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

func UsernameValidator(validate *validator.Validate) {
	validate.RegisterValidation("usernamevalidator", func(fl validator.FieldLevel) bool {
		usernameRegex := `^[a-zA-Z\d]{5,12}$`
		return regexp.MustCompile(usernameRegex).MatchString(fl.Field().String())
	})
}

func PasswordValidator(validate *validator.Validate) {
	validate.RegisterValidation("passwordvalidator", func(fl validator.FieldLevel) bool {
		passwordRegex := `^[a-zA-Z\d@_-]{8,20}$`
		password := fl.Field().String()
		ok := regexp.MustCompile(passwordRegex).MatchString(password)
		if !ok {
			return false
		}

		isSpesialCharacter := strings.ContainsAny(password, "@ | _ | -")

		isUpper := false
		isLower := false
		isNumber := false
	isPasswordLoop:
		for _, value := range password {
			if unicode.IsUpper(value) && unicode.IsLetter(value) && !isUpper {
				isUpper = true
			} else if unicode.IsLower(value) && unicode.IsLetter(value) && !isLower {
				isLower = true
			} else if _, err := strconv.Atoi(string(value)); err == nil {
				isNumber = true
			}

			if isUpper && isLower && isNumber {
				break isPasswordLoop
			}
		}

		if !isSpesialCharacter || !isUpper || !isLower || !isNumber {
			return false
		}
		return true
	})
}

func TelephoneValidator(validate *validator.Validate) {
	validate.RegisterValidation("telephonevalidator", func(fl validator.FieldLevel) bool {
		regexString := `^[\d+]{14}$`
		return regexp.MustCompile(regexString).MatchString(fl.Field().String())
	})
}

func SetValidator() (validate *validator.Validate) {
	validate = validator.New()
	UsernameValidator(validate)
	PasswordValidator(validate)
	TelephoneValidator(validate)
	return
}
