package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	LogFile = os.Getenv("HOME") + "/stream-alert-bot/logs/stream-alert-bot.log"
	logDir  = os.Getenv("HOME") + "/stream-alert-bot/logs"
)

var Log = logrus.New()

// InitLogger initializes the package logger.
// It removes any existing log file at LogFile, ensures the log directory (logDir)
// exists (0755), and opens or creates the log file with permissions 0644 in
// append mode. The package logger's output is set to both stdout and the log
// file. The logger is configured with a text formatter that uses full
// timestamps formatted as "2006-01-02 15:04:05" and forces colored output.
// The log level is set to Debug. If opening the log file fails, the error is
// logged and the process exits.
func InitLogger() {
	os.Remove(LogFile)
	_ = os.MkdirAll(logDir, 0755)

	file, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		Log.Errorf("failed to open log file: %v", err)
		os.Exit(0)
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, file))

	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
	})

	Log.SetLevel(logrus.DebugLevel)
}
