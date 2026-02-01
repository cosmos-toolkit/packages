// Package config fornece carregamento de env (env, dotenv), validação e defaults.
// Serve para CLI, cron, API e worker.
package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Loader carrega variáveis de ambiente. Load() pode ser chamado antes de Parse.
type Loader struct {
	loaded bool
}

// LoadFromFiles carrega arquivos .env (ex.: .env, .env.local).
// Se nenhum arquivo for passado, não faz nada. Erros de arquivo inexistente são ignorados.
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

// LoadFromFiles é um helper que carrega .env e retorna erro.
func LoadFromFiles(filenames ...string) error {
	var loader Loader
	return loader.LoadFromFiles(filenames...)
}

// Get retorna a variável de ambiente ou o valor default.
func Get(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// MustGet retorna a variável de ambiente ou panics se vazia.
func MustGet(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	panic("config: missing required env " + key)
}
