package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

var dataDir string

// InitData находит корень проекта и создаёт папку data
func InitData() {
	if dataDir != "" {
		return // уже инициализировано
	}

	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			log.Fatal("Не удалось найти go.mod — неизвестен корень проекта")
		}
		dir = parent
	}

	rootDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	dataDir = filepath.Join(rootDir, "data")
	if err := os.MkdirAll(dataDir, 0770); err != nil {
		log.Fatal(err)
	}
	err = godotenv.Load(filepath.Join(dataDir, ".env"))
	if err != nil {
		log.Fatal(err)
	}
}

// GetDataPath возвращает абсолютный путь к файлу в папке data
func GetDataPath(filename string) string {
	if dataDir == "" {
		InitData()
	}
	return filepath.Join(dataDir, filename)
}
