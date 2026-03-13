// Package logger holds logger implementation
package logger

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/w1lam/mc-pacman/internal/core/events"
)

type logger struct {
	logDirPath string
	buffer     []events.Event
	bufferSize int
	index      int
	count      int
	mu         sync.Mutex
}

const defaultBufferSize = 100

// New creates a new logger
func New(dirPath string) *logger {
	return &logger{
		logDirPath: dirPath,
		buffer:     make([]events.Event, defaultBufferSize),
		bufferSize: defaultBufferSize,
	}
}

// Log logs an event to log buffer
func (l *logger) Log(e events.Event) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.buffer[l.index] = e
	l.index = (l.index + 1) % l.bufferSize
	l.count++
}

// Close closes the logger, writing log to file if there are errors
func (l *logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.hasErrors() {
		return nil
	}

	return l.writeBufferToFile()
}

// HELPERS

func (l *logger) hasErrors() bool {
	for i := 0; i < l.bufferSize && i < l.count; i++ {
		if l.buffer[i].ErrorLvl == events.ERROR ||
			l.buffer[i].ErrorLvl == events.FATAL {
			return true
		}
	}

	return false
}

func (l *logger) getBufferInOrder() []events.Event {
	if l.count < l.bufferSize {
		return l.buffer[:l.count]
	}

	result := make([]events.Event, l.bufferSize)
	for i := 0; i < l.bufferSize; i++ {
		result[i] = l.buffer[(l.index+i)%l.bufferSize]
	}

	return result
}

func (l *logger) writeBufferToFile() error {
	dir := filepath.Dir(l.logDirPath)
	if err := os.MkdirAll(dir, 0o766); err != nil {
		return err
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	logPath := filepath.Join(dir, timestamp+".log")

	f, err := os.Create(logPath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	defer w.Flush()

	fmt.Fprintf(w, "=== MC-PACMAN ERROR LOG ===")
	fmt.Fprintf(w, "Generated: %s\n\n", time.Now().Format("2006-01-01_15:04:05"))

	events := l.getBufferInOrder()
	for _, e := range events {
		l.formatEvent(w, e)
	}
	return nil
}
