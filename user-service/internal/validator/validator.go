package validator

import (
	"fmt"
	"regexp"
	"strings"

	validator_v10 "github.com/go-playground/validator/v10"
)

var validate *validator_v10.Validate

var validUsername = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-_.]{1,20}(@[a-zA-Z][a-zA-Z0-9-_.]{1,20})?$`)

func InitValidator() {
	validate = validator_v10.New(validator_v10.WithRequiredStructEnabled())

	_ = validate.RegisterValidation("not_whitespace", notWhitespace)
	_ = validate.RegisterValidation("uri", validateURI)
	_ = validate.RegisterValidation("username", validateUsername)
}

func Validate(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		return fmt.Errorf("validator: %w", err)
	}
	return nil
}

func ValidateUsername(username string) error {
	// Проверка на соответствие регулярному выражению
	if !validUsername.MatchString(username) {
		return fmt.Errorf("validator: invalid username")
	}

	return nil
}

func notWhitespace(fl validator_v10.FieldLevel) bool {
	field := fl.Field().String()
	if field == "" {
		return true
	}

	trimmed := strings.TrimSpace(field)

	return trimmed != ""
}

func validateURI(fl validator_v10.FieldLevel) bool {
	uri := fl.Field().String()
	if strings.HasPrefix(uri, "http://") {
		return len(uri) > len("https://")
	}
	if strings.HasPrefix(uri, "https://") {
		return len(uri) > len("https://")
	}

	return false
}

func validateUsername(fl validator_v10.FieldLevel) bool {
	field := fl.Field().String()

	return validUsername.MatchString(field)
}
