package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
)

type LogEntry struct {
	Level   string
	Message string
}

type Logger struct {
	logCh    chan *LogEntry
	writer   io.Writer
	wg       sync.WaitGroup
	stop     chan struct{}
	stopOnce sync.Once
}

func New(writer io.Writer) *Logger {
	return &Logger{
		logCh:  make(chan *LogEntry),
		writer: writer,
		stop:   make(chan struct{}),
	}
}

func (l *Logger) Run() {
	go func() {
		defer l.wg.Done()
		for {
			select {
			case entry := <-l.logCh:
				// Maybe printing out all messages could be optional
				// or possibly only print out certain levels, could do?
				//fmt.Printf("[%s] %s\r\n", entry.Level, entry.Message)
				fmt.Fprintf(l.writer, "[%s] %s\r\n", entry.Level, entry.Message)
				l.wg.Add(1)
			case <-l.stop:
				return
			}
		}
	}()
}

func (l *Logger) Log(level, message string, any ...any) {
	l.wg.Add(1)
	defer l.wg.Done()
	select {
	case <-l.stop:
	// Don't send message if the logger is stopped
	default:
		l.logCh <- &LogEntry{Level: level, Message: fmt.Sprintf(message, any...)}
	}
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
	file, ok := l.writer.(*os.File)
	if ok {
		file.Close()
	}
}
