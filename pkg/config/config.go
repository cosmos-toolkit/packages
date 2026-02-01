// Package config provides env loading (env, dotenv), validation and defaults.
// Used by CLI, cron, API and worker.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Loader loads environment variables. Load() can be called before Parse.
type Loader struct {
	loaded bool
}

// LoadFromFiles loads .env files (e.g. .env, .env.local).
// If no files are passed, does nothing. Missing file errors are ignored.
func (l *Loader) LoadFromFiles(filenames ...string) error {
	if len(filenames) == 0 {
		return nil
	}
	if err := godotenv.Load(filenames...); err != nil && !os.IsNotExist(err) {
		return err
	}
	l.loaded = true
	return nil
}

// LoadFromFiles is a helper that loads .env and returns an error.
func LoadFromFiles(filenames ...string) error {
	var loader Loader
	return loader.LoadFromFiles(filenames...)
}

// Get returns the environment variable or the default value.
func Get(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// MustGet returns the environment variable or panics if empty.
func MustGet(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("config: missing required env " + key)
}
