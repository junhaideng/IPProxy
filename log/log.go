package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var levelMap = map[string]logrus.Level{
	"info":    logrus.InfoLevel,
	"debug":   logrus.DebugLevel,
	"trace":   logrus.TraceLevel,
	"warning": logrus.WarnLevel,
	"error":   logrus.ErrorLevel,
}

const (
	ConsoleMode = "console"
	FileMode    = "file"
)

func Init() {
	fmt.Println("initialize log")
	level := viper.GetString("log.level")
	logrus.SetLevel(levelMap[level])

	mode := viper.GetString("log.mode")

	if strings.EqualFold(mode, ConsoleMode) {

	} else if strings.EqualFold(mode, FileMode) {
		filename := viper.GetString("log.filename")
		f, err := os.Create(filename)
		if err != nil {
			panic(err)
		}
		logrus.SetOutput(f)
	} else {
		logrus.Info("log mode specified not legal, use console mode default")
	}
}
