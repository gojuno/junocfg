# go.log [![GoDoc](https://godoc.org/github.com/mguzelevich/go.log?status.svg)](http://godoc.org/github.com/mguzelevich/go.log) [![Build Status](https://travis-ci.org/mguzelevich/go.log.svg?branch=master)](https://travis-ci.org/mguzelevich/go.log)


go logger wrapper


# usage

NOTE: use `nil` instead `ioutil.Discard` for disable log level

init:

```
import (
	"os"

	"github.com/mguzelevich/go.log"
)

func main() {
	log.InitLoggers(&log.Logger{
		nil,
		nil,
		os.Stdout,
		os.Stdout,
		os.Stderr,
	})

	log.InitLoggers(&log.Logger{ Error: os.Stderr })

}
```

usage

```
import "github.com/mguzelevich/go-log"

log.Debug.Printf("some debug message")
...

log.Error.Printf("some error message %s", err)

```
*/
