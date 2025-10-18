package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"main.go/internal/util"
	"os"
)

const LogFile = "/var/log/stream-alert-bot/stream-alert-bot.log"
const logDir = "/var/log/stream-alert-bot"

var Log = logrus.New()

func InitLogger() {
	os.Remove(LogFile)
	_ = os.MkdirAll(logDir, 0755)

	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Log.Errorf("failed to open log file: %v", err)
		util.WaitForSignal()
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, file))

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	Log.SetLevel(logrus.DebugLevel)
}
