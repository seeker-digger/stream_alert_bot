package config

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	l "main.go/internal/logger"
)

var envVarSlice = []string{"KICK_CLIENT_ID", "KICK_CLIENT_SECRET", "TELEGRAM_BOT_API"}

const envFileText = "" +
	"KICK_CLIENT_ID=\"\"\n" +
	"KICK_CLIENT_SECRET=\"\"\n" +
	"TELEGRAM_BOT_API=\"\""

func Init() {
	home, err := os.UserHomeDir()
	if err != nil {
		l.Log.Panic(err)
	}
	environmentFile := path.Join(home, "/stream-alert-bot/.env")

	l.InitLogger()

	err = godotenv.Load(environmentFile)
	if err != nil {
		if strings.Contains(err.Error(), ".env: no such file or directory") {
			err = os.MkdirAll(filepath.Dir(environmentFile), 0666)
			if err != nil {
				log.Panic(err)
			}
			err = os.WriteFile(environmentFile, []byte(envFileText), 0666)
			if err != nil {
				log.Panic(err)
			}
			log.Println("Please fill the .env file on this way: " + environmentFile)
		} else {
			log.Fatal(err)
		}
	}
	for _, i := range envVarSlice {
		a := os.Getenv(i)
		if a == "" {
			l.Log.Error("Please set " + i + " and others else in the .env file on this way: " + environmentFile)
			os.Exit(0)
		}
	}
}

func GetDataPath(filename string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		l.Log.Panic(err)
	}
	workingDirectory := path.Join(home, "/stream-alert-bot")
	return filepath.Join(workingDirectory, filename)
}
