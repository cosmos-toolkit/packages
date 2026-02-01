// Package validator fornece um wrapper do go-playground/validator/v10:
// validação de structs, mensagens padronizadas e reuso entre API e CLI.
package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// V é o validador compartilhado.
var V = validator.New()

// Validate valida a struct v e retorna erro com mensagens amigáveis.
func Validate(v interface{}) error {
	if err := V.Struct(v); err != nil {
		return Translate(err)
	}
	return nil
}

// Translate converte erros do validator em mensagens padronizadas.
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
	return fmt.Errorf("validação: %s", strings.Join(msgs, "; "))
}

func fieldError(e validator.FieldError) string {
	field := e.Field()
	// snake_case para API
	name := toSnake(field)
	switch e.Tag() {
	case "required":
		return name + " é obrigatório"
	case "email":
		return name + " deve ser um e-mail válido"
	case "min":
		return name + " deve ter no mínimo " + e.Param() + " caracteres"
	case "max":
		return name + " deve ter no máximo " + e.Param() + " caracteres"
	case "oneof":
		return name + " deve ser um de: " + e.Param()
	default:
		return name + " é inválido (" + e.Tag() + ")"
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
