package logger

import (
	"io/fs"
	"log"
	"os"
)

var isLoggerEnabled bool

func SetLoggerFile(path string) (*os.File, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		// create file
		logFile, err := os.Create(path)
		if err != nil {
			return nil, err
		}

		log.SetOutput(logFile)
		isLoggerEnabled = true
		return logFile, nil
	} else if err != nil {
		return nil, err
	} else {
		logFile, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, fs.ModeAppend)
		if err != nil {
			return nil, err
		}
		log.SetOutput(logFile)
		isLoggerEnabled = true
		return logFile, nil
	}
}

func Println(messages ...any) {
	if isLoggerEnabled {
		log.Println(messages...)
	}
}
