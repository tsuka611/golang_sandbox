package log

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

type LogLevel int8
type appWriter struct {
	mu sync.Mutex
	w  io.Writer
}

func (l *appWriter) setWriter(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.w = w
}
func (l *appWriter) Write(p []byte) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.w.Write(p)
}

const (
	L_TRACE = LogLevel(1) << iota
	L_DEBUG = LogLevel(1) << iota
	L_INFO  = LogLevel(1) << iota
	L_WARN  = LogLevel(1) << iota
	L_ERROR = LogLevel(1) << iota
)

var (
	wTrace = &appWriter{w: ioutil.Discard}
	wDebug = &appWriter{w: ioutil.Discard}
	wInfo  = &appWriter{w: ioutil.Discard}
	wWarn  = &appWriter{w: ioutil.Discard}
	wError = &appWriter{w: ioutil.Discard}

	TRACE = log.New(wTrace, "TRACE ", log.Ldate|log.Ltime|log.Lshortfile)
	DEBUG = log.New(wDebug, "DEBUG ", log.Ldate|log.Ltime|log.Lshortfile)
	INFO  = log.New(wInfo, "INFO ", log.Ldate|log.Ltime|log.Lshortfile)
	WARN  = log.New(wWarn, "WARN ", log.Ldate|log.Ltime|log.Lshortfile)
	ERROR = log.New(wError, "ERROR ", log.Ldate|log.Ltime|log.Lshortfile)

	outputs = map[LogLevel]*appWriter{
		L_TRACE: wTrace,
		L_DEBUG: wDebug,
		L_INFO:  wInfo,
		L_WARN:  wWarn,
		L_ERROR: wError,
	}
)

func SetLogLevel(level LogLevel) {
	for k, v := range outputs {
		if k < level {
			v.setWriter(ioutil.Discard)
		} else {
			v.setWriter(os.Stderr)
		}
	}
}

func init() {
	log.SetFlags(INFO.Flags())
	log.SetOutput(os.Stderr)
	SetLogLevel(L_INFO)
}
