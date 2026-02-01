package config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Validator validates structs with "env", "required" or "default" tags.
// Minimal usage: struct with string fields and tag `config:"ENV_KEY,required"` or `config:"ENV_KEY,default=value"`.
const tagName = "config"

// Validate fills defaults and validates required fields in v (struct pointer).
// Supported tags: config:"ENV_KEY,required" and config:"ENV_KEY,default=value".
// If a field has "required", the value must be non-empty (after filling from env or default).
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("config: Validate expects *struct, got %T", v)
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
			return fmt.Errorf("config: %s (env %s) is required", f.Name, envKey)
		}
	}
	return nil
}
