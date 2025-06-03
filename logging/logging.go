package logging

import (
	"log"
	"os"
	"sync"
)

type FLogger struct {
	logFile *os.File
	mu      sync.Mutex
}

var DefaultLogger *FLogger
var TranslationsLogger *FLogger

func init() {
	var err error
	DefaultLogger, err = NewFLogger("logging/app.log")
	if err != nil {
		log.Printf("Failed to initialize default logger: %v", err)
		return
	}

	TranslationsLogger, err = NewFLogger("logging/translations.log")
	if err != nil {
		log.Printf("Failed to initialize translations logger: %v", err)
		return
	}
}

func NewFLogger(filename string) (*FLogger, error) {
	logFile, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &FLogger{logFile: logFile}, nil
}

func (l *FLogger) LogError(err error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.SetOutput(l.logFile)
	log.Println(err)
}

func (l *FLogger) LogErrorF(format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	log.SetOutput(l.logFile)
	log.Printf(format, args...)
}
