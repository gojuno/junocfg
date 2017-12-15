package log

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var timestamp string

func initAndLog(logger *Logger) {
	InitLoggers(logger)

	Stdout.Printf("Stdout.Printf")
	Stderr.Printf("Stderr.Printf")
	Trace.Printf("Trace.Printf")
	Debug.Printf("Debug.Printf")
	Info.Printf("Info.Printf")
	Warning.Printf("Warning.Printf")
	Error.Printf("Error.Printf")
}

var tests = map[string]struct {
	logger  *Logger
	enabled []string
}{
	"empty_init": {
		logger:  nil,
		enabled: []string{"stdout", "stderr"},
	},
	"trace": {
		logger:  &Logger{Trace: os.Stdout},
		enabled: []string{"stdout", "stderr", "trace"},
	},
	"all": {
		logger:  &Logger{os.Stdout, os.Stdout, os.Stdout, os.Stdout, os.Stderr},
		enabled: []string{"stdout", "stderr", "trace", "debug", "info", "warning", "error"},
	},
}

func getExpected(ts string, loggers ...string) string {
	baseLine := 21
	out := []string{}
	for _, l := range loggers {
		switch l {
		case "stdout":
			out = append(out, fmt.Sprintf("%s log_test.go:%d: Stdout.Printf", ts, baseLine))
		case "stderr":
			out = append(out, fmt.Sprintf("%s log_test.go:%d: Stderr.Printf", ts, baseLine+1))
		case "trace":
			out = append(out, fmt.Sprintf("TRACE: %s log_test.go:%d: Trace.Printf", ts, baseLine+2))
		case "debug":
			out = append(out, fmt.Sprintf("DEBUG: %s log_test.go:%d: Debug.Printf", ts, baseLine+3))
		case "info":
			out = append(out, fmt.Sprintf("INFO: %s Info.Printf", ts)) // baseLine+4
		case "warning":
			out = append(out, fmt.Sprintf("WARNING: %s log_test.go:%d: Warning.Printf", ts, baseLine+5))
		case "error":
			out = append(out, fmt.Sprintf("ERROR: %s log_test.go:%d: Error.Printf", ts, baseLine+6))
		default:
			panic("unknown")
		}
	}
	out = append(out, "")
	return strings.Join(out, "\n")
}

func _TestRun(t *testing.T) {
	t.Logf("run %v %v", loggers["trace"], *loggers["trace"].logger)
	InitLoggers(&Logger{
		Trace: os.Stdout,
	})

	t.Logf("run %v %v", loggers["trace"], *loggers["trace"].logger)
	t.Logf("run %v", Levels())

	assert := assert.New(t)
	assert.NotNil(nil)
}

func TestLog(t *testing.T) {
	// https://medium.com/agrea-technogies/basic-testing-patterns-in-go-d8501e360197
	name := os.Getenv("GO_TEST_FORKED_PROCESS")
	tst := tests[name]
	if name != "" {
		// fmt.Printf("run %s %v %v\n", name, tst, Levels())
		initAndLog(tst.logger)
		os.Exit(0)
		return
	}

	assert := assert.New(t)
	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		out, err := helperCommand(t, testName)
		assert.Nil(err)
		// assert.Equal(test.enabled, Levels(), "")
		assert.Equal(getExpected(timestamp, test.enabled...), string(out), "")
	}
}

func ExampleLog_stdout() {
	// init logger
	InitLoggers(&Logger{
		ioutil.Discard,
		ioutil.Discard,
		os.Stdout,
		os.Stdout,
		os.Stderr,
	})

	Stdout.Printf("Stdout.Printf")
	Stderr.Printf("Stderr.Printf")
}

func helperCommand(t *testing.T, testName string) (string, error) {
	cs := []string{fmt.Sprintf("-test.run=%s", helperCaller()), "--"}
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{fmt.Sprintf("GO_TEST_FORKED_PROCESS=%s", testName)}

	out, err := cmd.CombinedOutput()

	return string(out), err
}

func helperCaller() string {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	me := runtime.FuncForPC(pc)
	if me == nil {
		return "unnamed"
	}
	full := me.Name()
	p := strings.Split(full, ".")
	return p[len(p)-1]
}

func init() {
	timestamp = time.Now().UTC().Format("2006/01/02 15:04:05")
}
