package config

import (
	"errors"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
)

var errNoSuchFile = errors.New("open /etc/alert-bot/.env: no such file or directory")

var envVarSlice = []string{"KICK_CLIENT_ID", "KICK_CLIENT_SECRET", "TELEGRAM_BOT_API"}

const envFileText = "" +
	"KICK_CLIENT_ID=\"\"\n" +
	"KICK_CLIENT_SECRET=\"\"\n" +
	"TELEGRAM_BOT_API=\"\""

const environmentFile = "/etc/alert-bot/.env"
const workingDirectory = "/var/lib/alert-bot"

func InitData() {
	err := godotenv.Load(environmentFile)
	if err != nil {
		if errors.Is(err, errNoSuchFile) {
			file, err := os.Create(environmentFile)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			_, err = file.Write([]byte(envFileText))
			if err != nil {
				log.Fatal(err)
			}
			log.Fatal("Please fill the .env file on this way: " + environmentFile)
		} else {
			log.Fatal(err)
		}
	}
	for _, i := range envVarSlice {
		a := os.Getenv(i)
		if a == "" {
			log.Fatal("Please set " + i + " in the .env file on this way: " + environmentFile)
		}
	}
}

func GetDataPath(filename string) string {
	return filepath.Join(workingDirectory, filename)
}
