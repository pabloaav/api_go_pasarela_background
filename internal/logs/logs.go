package logs

import (
	"log"
	"os"
)

var (
	infoLogger  = log.New(os.Stdout, "", 0)
	errorLogger = log.New(os.Stdout, "", log.Lshortfile|log.Ltime)
)

///Info ...Log Info
func Info(s interface{}) {
	infoLogger.Print(s)
}

///Error ...Log Errors
func Error(s interface{}) {
	errorLogger.Print(s)
}
