package logs

import (
	//Inbuild packages

	"io"
	"os"

	//Third-party packages
	"github.com/sirupsen/logrus"
	call "github.com/t-tomalak/logrus-easy-formatter"
)

var Log *logrus.Logger

func LoggerConfiguration() {
	file, err := os.OpenFile("logs.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}

	Log = &logrus.Logger{
		Out:   io.MultiWriter(file, os.Stdout),
		Level: logrus.DebugLevel,
		Formatter: &call.Formatter{
			TimestampFormat: "02-01-2006 15:04:05",
			LogFormat:       "[%lvl%]: %time% - %msg%\n",
		},
	}
}
