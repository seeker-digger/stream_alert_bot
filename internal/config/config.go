package config

import (
	"github.com/joho/godotenv"
	"log"
	l "main.go/internal/logger"
	"main.go/internal/util"
	"os"
	"path/filepath"
	"strings"
)

var envVarSlice = []string{"KICK_CLIENT_ID", "KICK_CLIENT_SECRET", "TELEGRAM_BOT_API"}

const envFileText = "" +
	"KICK_CLIENT_ID=\"\"\n" +
	"KICK_CLIENT_SECRET=\"\"\n" +
	"TELEGRAM_BOT_API=\"\""

const environmentFile = "/etc/alert-bot/.env"
const workingDirectory = "/var/lib/alert-bot"

func Init() {
	l.InitLogger()

	err := godotenv.Load(environmentFile)
	if err != nil {
		if strings.Contains(err.Error(), ".env: no such file or directory") {
			err = os.MkdirAll(filepath.Dir(environmentFile), 0766)
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
			util.WaitForSignal()
		}
	}
}

func GetDataPath(filename string) string {
	return filepath.Join(workingDirectory, filename)
}
