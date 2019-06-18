package ini

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	u "utxo/utils"
)

func initLog() {
	file, err := os.OpenFile(u.GetProjectPath("utxo")+"utxo.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	multiWrite := io.MultiWriter(os.Stdout, file)
	if err == nil {
		log.SetOutput(multiWrite)
	} else {
		log.Info("Failed to log to file")
		os.Exit(3)
	}
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
}

//GetLog is get logrus obj
func GetLog() *logrus.Logger {
	return log
}
