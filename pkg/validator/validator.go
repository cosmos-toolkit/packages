// Package validator provides a wrapper for go-playground/validator/v10:
// struct validation, standardized messages, and reuse between API and CLI.
package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// V is the shared validator instance.
var V = validator.New()

// Validate validates struct v and returns an error with user-friendly messages.
func Validate(v interface{}) error {
	if err := V.Struct(v); err != nil {
		return Translate(err)
	}
	return nil
}

// Translate converts validator errors into standardized messages.
func Translate(err error) error {
	if err == nil {
		return nil
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return err
	}
	var msgs []string
	for _, e := range errs {
		msgs = append(msgs, fieldError(e))
	}
	return fmt.Errorf("validation: %s", strings.Join(msgs, "; "))
}

func fieldError(e validator.FieldError) string {
	field := e.Field()
	// snake_case for API
	name := toSnake(field)
	switch e.Tag() {
	case "required":
		return name + " is required"
	case "email":
		return name + " must be a valid email"
	case "min":
		return name + " must have at least " + e.Param() + " characters"
	case "max":
		return name + " must have at most " + e.Param() + " characters"
	case "oneof":
		return name + " must be one of: " + e.Param()
	default:
		return name + " is invalid (" + e.Tag() + ")"
	}
}

func toSnake(s string) string {
	var b strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			b.WriteByte('_')
		}
		if r >= 'A' && r <= 'Z' {
			r = r + 32
		}
		b.WriteRune(r)
	}
	return b.String()
}
