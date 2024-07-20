package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARNING"
	LogLevelError = "ERROR"
	LogLevelNone  = "NONE"
)

type LogEntry struct {
	Level   string
	Message string
}

type Logger struct {
	logCh      chan *LogEntry
	writer     io.Writer
	wg         sync.WaitGroup
	stop       chan struct{}
	stopOnce   sync.Once
	debugLevel string // Flag for enabling debug logging
}

func New(writer io.Writer, debugLevel string) *Logger {
	return &Logger{
		logCh:      make(chan *LogEntry),
		writer:     writer,
		stop:       make(chan struct{}),
		debugLevel: debugLevel,
	}
}

func (l *Logger) Run() {
	go func() {
		for {
			select {
			case entry := <-l.logCh:
				fmt.Fprintf(l.writer, "[%s] %s\r\n", entry.Level, entry.Message)
				if l.shouldPrintDebug(entry.Level) {
					fmt.Printf("[%s] %s\r\n", entry.Level, entry.Message)
				}
				l.wg.Done()
			case <-l.stop:
				close(l.logCh) // Just for explicitness
				return
			}
		}
	}()
}

func (l *Logger) Log(level, format string, any ...any) {
	l.wg.Add(1)
	select {
	case <-l.stop:
	// Don't send message if the logger is stopped
	default:
		l.logCh <- &LogEntry{Level: level, Message: fmt.Sprintf(format, any...)}
	}
}

func (l *Logger) shouldPrintDebug(level string) bool {

	return l.debugLevel == "INFO" || (l.debugLevel == LogLevelWarn && (level == LogLevelWarn || level == LogLevelError)) || level == l.debugLevel
}

// Stop signals the logger to stop processing messages
func (l *Logger) Stop() {
	l.stopOnce.Do(func() {
		close(l.stop)
	})
}

func (l *Logger) Wait() {
	l.Stop()
	l.wg.Wait()
	// If the writer is a file, then close it as we are done with it.
	if file, ok := l.writer.(*os.File); ok {
		file.Close()
	}
}
