package log

import (
	"crypto/rand"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

type loggerDesc struct {
	logger  **log.Logger
	enabled bool
}

var (
	Trace   *log.Logger
	Debug   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger

	Stdout *log.Logger
	Stderr *log.Logger

	levels  []string
	loggers map[string]*loggerDesc
)

type Logger struct {
	Trace   io.Writer
	Debug   io.Writer
	Info    io.Writer
	Warning io.Writer
	Error   io.Writer
}

func InitLoggers(logger *Logger) {
	if logger == nil {
		return
	}
	set("trace", logger.Trace, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	set("debug", logger.Debug, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	set("info", logger.Info, "INFO: ", log.Ldate|log.Ltime|log.LUTC)
	set("warning", logger.Warning, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	set("error", logger.Error, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)

	// Stdout = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
	// Stderr = log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)
}

func set(name string, out io.Writer, prefix string, flag int) error {
	desc, ok := loggers[name]
	desc.enabled = out != nil
	if out == nil {
		out = ioutil.Discard
	}

	if ok {
		*(desc.logger) = log.New(out, prefix, flag)
	} else {
		return fmt.Errorf("unknown logger level %s", name)
	}
	return nil
}

func init() {
	levels = []string{"stdout", "stderr", "trace", "debug", "info", "warning", "error"}
	loggers = map[string]*loggerDesc{
		"stdout":  &loggerDesc{&Stdout, false},
		"stderr":  &loggerDesc{&Stderr, false},
		"trace":   &loggerDesc{&Trace, false},
		"debug":   &loggerDesc{&Debug, false},
		"info":    &loggerDesc{&Info, false},
		"warning": &loggerDesc{&Warning, false},
		"error":   &loggerDesc{&Error, false},
	}

	set("stdout", os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	set("stderr", os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile|log.LUTC)
	for _, n := range []string{"trace", "debug", "info", "warning", "error"} {
		set(n, nil, "", log.Ldate)
	}
}

func Levels() []string {
	out := []string{}
	for _, n := range levels {
		if loggers[n].enabled {
			out = append(out, n)
		}
	}
	return out
}

// UUID generates a random UUID according to RFC 4122
func UUID() (string, error) {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return "", err
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:]), nil
}
