package core

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// AppError represents an application error
type AppError struct {
	Time    time.Time         `json:"time"`
	Source  string            `json:"source"`
	Message string            `json:"message"`
	Fatal   bool              `json:"fatal"`
	Context map[string]string `json:"context,omitempty"`

	err error `json:"-"`
}

var (
	errCh   = make(chan AppError, 256)
	once    sync.Once
	logFile *os.File
)

// Start starts the error logging system
func Start(logPath string) error {
	var err error
	once.Do(func() {
		logFile, err = os.OpenFile(
			logPath,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0o644,
		)
		if err != nil {
			return
		}

		go processErrors()
	})
	return err
}

func processErrors() {
	writer := bufio.NewWriter(logFile)
	defer writer.Flush()

	enc := json.NewEncoder(writer)

	for e := range errCh {
		_ = enc.Encode(e)

		if e.Fatal {
			fmt.Println("FATAL:", e.Message)
			writer.Flush()
			os.Exit(1)
		}
	}
}

// Report reports an error to the error logging system
func Report(source string, err error) {
	if err == nil {
		return
	}

	select {
	case errCh <- AppError{
		Time:    time.Now(),
		Source:  source,
		Message: err.Error(),
		err:     err,
	}:
	default:
		fmt.Println("ERROR:", err)
	}
}

// ReportCtx reports an error with context to the error logging system
func ReportCtx(source string, err error, ctx map[string]string) {
	if err == nil {
		return
	}

	select {
	case errCh <- AppError{
		Time:    time.Now(),
		Source:  source,
		Message: err.Error(),
		err:     err,
		Context: ctx,
	}:
	default:
		fmt.Println("ERROR:", err)
	}
}
