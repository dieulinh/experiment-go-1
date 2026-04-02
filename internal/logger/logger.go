package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	// Use JSON format (good for production)
	Log.SetFormatter(&logrus.JSONFormatter{})
	// Log.SetReportCaller(True)

	// Output to console
	Log.SetOutput(os.Stdout)

	// Set log level
	Log.SetLevel(logrus.DebugLevel)
}
