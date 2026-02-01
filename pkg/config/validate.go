package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Validator valida structs com tags "env" e "required" ou "default".
// Uso mínimo: struct com campos string e tag `config:"ENV_KEY,required"` ou `config:"ENV_KEY,default=valor"`.
const tagName = "config"

// Validate preenche defaults e valida campos obrigatórios em v (struct pointer).
// Tags suportadas: config:"ENV_KEY,required" e config:"ENV_KEY,default=valor".
// Se um campo tem "required", o valor deve ser não vazio (após preencher do env ou default).
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config: Validate espera *struct, obteve %T", v)
	}
	val = val.Elem()
	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.CanSet() {
			continue
		}
		f := typ.Field(i)
		tag := f.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		parts := strings.Split(tag, ",")
		envKey := strings.TrimSpace(parts[0])
		var required bool
		var defaultVal string
		for _, p := range parts[1:] {
			p = strings.TrimSpace(p)
			if p == "required" {
				required = true
			} else if strings.HasPrefix(p, "default=") {
				defaultVal = strings.TrimSpace(strings.TrimPrefix(p, "default="))
			}
		}
		if field.Kind() != reflect.String {
			continue
		}
		current := field.String()
		if current == "" {
			if envVal := os.Getenv(envKey); envVal != "" {
				field.SetString(envVal)
				current = envVal
			} else if defaultVal != "" {
				field.SetString(defaultVal)
				current = defaultVal
			}
		}
		if required && current == "" {
			return fmt.Errorf("config: %s (env %s) é obrigatório", f.Name, envKey)
		}
	}
	return nil
}
