package logstream

import (
	"fmt"
	"io"
	"log"
)

type Log struct {
	log *log.Logger
}

func Logger(out io.Writer) (l *Log) {
	return &Log{log: log.New(out, "[logstream] ", log.Ldate|log.Ltime|log.Lmicroseconds)}
}

func (l *Log) Logf(f string, args ...interface{}) {
	l.log.Output(2, fmt.Sprintf(f, args...))
}
