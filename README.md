# glog

A log designed for critical mission.

The glog package provides a fast and simple logger dedicated to TEXT/JSON output

## Features
* Sampling and fast
* Low to zero allocation
* Level logging
* Additional fields
* `context.Context` integration
* JSON and TEXT encoding formats
* User-define Encoder
* User-define Executor

## Installation

```bash
go get -u github.com/DataWorkbench/glog
```
Used in go modules
```bash
go get -insecure github.com/DataWorkbench/glog
```

## Quick Start

#### Simple Logging Example
```go
package main

import (
	"github.com/DataWorkbench/glog"
)

func main() {
    // By default, the log message will output to os.Stdout
	l := glog.NewDefault()

	l.Debug().Msg("HelloWorld").String("s1", "v1").Int64("i1", 1).Fire()

	/* Output:
	2020-11-04T18:01:40.945196+08:00 [debug] HelloWorld s1=v1 i1=1
	*/
}
```

#### Add fixed fields into logger
```go
package main

import (
	"github.com/DataWorkbench/glog"
)

func main() {
	l := glog.NewDefault()
	l.WithFields().AddString("requestid", "8da3aceea1ba")

	l.Debug().Msg("HelloWorld").String("s1", "v1").Int64("i1", 1).Fire()

	/* Output:
	2020-11-04T18:00:56.050655+08:00 [debug] HelloWorld s1=v1 i1=1 requestid=8da3aceea1ba
	*/
}
```

#### Set the logger level
```go
package main

import (
	"github.com/DataWorkbench/glog"
)

func main() {
	l := glog.NewDefault()
	l.WithLevel(glog.InfoLevel)
	
	// The debug log will be ignored
	l.Debug().Msg("Hello Debug Message").Fire()
	l.Info().Msg("Hello Info Message").Fire()

	/* Output:
	2020-11-04T17:58:23.958602+08:00 [info] Hello Info Message
	*/
}
```

#### Set Time Format
```go
package main

import (
	"time"

	"github.com/DataWorkbench/glog"
)

func main() {
	l := glog.NewDefault()
	l.WithTimeLayout(time.RFC822)

	l.Debug().Msg("HelloWorld").String("s1", "v1").Int64("i1", 1).Fire()

	/* Output:
	04 Nov 20 18:05 CST [debug] HelloWorld s1=v1 i1=1
	*/
}
``` 

#### Add caller info into log message
```go
package main

import (
	"github.com/DataWorkbench/glog"
)

func main() {
	l := glog.NewDefault()
	l.WithCaller(true)

	l.Debug().Msg("HelloWorld").String("s1", "v1").Int64("i1", 1).Fire()

	/* Output:
	2020-11-04T18:06:13.151354+08:00 [debug] HelloWorld s1=v1 i1=1 (github.com/DataWorkbench/glog/examples/main.go:11)
	*/
}
``` 

#### Write log message into file
```go
package main

import (
	"fmt"

	"github.com/DataWorkbench/glog"
)

func main() {
	logfile := "/tmp/testglog.log"
	executor, err := glog.FileExecutor(logfile, nil)
	if err != nil {
		fmt.Println("open file error:", err)
		return
	}
	defer executor.Close()

	l := glog.NewDefault()
	l.WithExecutor(executor)

	l.Debug().Msg("Hello logfile").Fire()

	/* Output:
	$ cat /tmp/testglog.log
	2020-11-04T18:22:21.726335+08:00 [debug] Hello logfile
	*/
}
```

#### Use with context.Context
```go
package main

import (
	"context"

	"github.com/DataWorkbench/glog"
)

func main() {
	// store logger into a context.Value
	ctx := glog.WithContext(context.Background(), glog.NewDefault())

	// get logger from context.Value
	l := glog.FromContext(ctx)

	l.Debug().Msg("Hello Context").Fire()

	/* Output:
	2020-11-04T21:15:21.002094+08:00 [debug] Hello Context
	*/
}
```

#### Use JSON Format
```go
package main

import (
	"github.com/DataWorkbench/glog"
)

func main() {
	l := glog.NewDefault()
	l.WithEncoderFunc(glog.JSONEncoder)
	l.WithFields().AddString("requestid", "8da3aceea1ba")

	l.Debug().Msg("HelloWorld").String("s1", "v1").Int64("i1", 1).Fire()

	/* Output:
	{"time":"2020-11-04T18:27:41.080215+08:00","level":"debug","message":"HelloWorld","s1":"v1","i1":1"requestid":"8da3aceea1ba"}
	*/
}
```

#### Clone from a exits logger
```go
package main

import (
	"github.com/DataWorkbench/glog"
)

func main() {
	l := glog.NewDefault()
	l.WithFields().AddString("requestid", "8da3aceea1ba")

	l.Debug().Msg("HelloWorld One").Fire()

	/* Output:
	2020-11-04T18:34:32.816702+08:00 [debug] HelloWorld One requestid=8da3aceea1ba
	*/

	// clone logger will
	lc := l.Clone()
	lc.WithFields().AddString("dup1", "dup-value")

	lc.Debug().Msg("Hello Clone Logger").Fire()

	/* Output:
	2020-11-04T18:34:32.816847+08:00 [debug] Hello Clone Logger requestid=8da3aceea1ba dup1=dup-value
	*/

	// any changed in close logger does not affects the sources
	lc.WithCaller(true)
	l.Debug().Msg("HelloWorld Two").Fire()

	/* Output:
	2020-11-04T18:34:32.816851+08:00 [debug] HelloWorld Two requestid=8da3aceea1ba
	*/
}
```

#### Write log message into multiple file
```go
package main

import (
	"fmt"

	"github.com/DataWorkbench/glog"
)

func main() {
	logFile := "/tmp/testglog.log"
	errLogFile := "/tmp/testglog.log.wf"

	e1, err := glog.FileExecutor(logFile, glog.MatchGELevel(glog.DebugLevel))
	if err != nil {
		fmt.Println("open log file fail:", err)
		return
	}

	e2, err := glog.FileExecutor(errLogFile, glog.MatchGELevel(glog.ErrorLevel))
	if err != nil {
		fmt.Println("open error log file fail:", err)
		return
	}

	l := glog.NewDefault()
	l.WithExecutor(glog.MultipleExecutor(e1, e2))

	l.Debug().Msg("DebugMessage").Fire()
	l.Info().Msg("InfoMessage").Fire()
	l.Error().Msg("ErrorMessage").Fire()
	l.Fatal().Msg("FatalMessage").Fire()

	/* Output:
	$cat /tmp/testglog.log
	2020-11-04T21:09:50.828122+08:00 [debug] DebugMessage
	2020-11-04T21:09:50.828354+08:00 [info] InfoMessage
	2020-11-04T21:09:50.828365+08:00 [error] ErrorMessage
	2020-11-04T21:09:50.828394+08:00 [fatal] FatalMessage

	$cat /tmp/testglog.log.wf
	2020-11-04T21:09:50.828365+08:00 [error] ErrorMessage
	2020-11-04T21:09:50.828394+08:00 [fatal] FatalMessage
	*/
}
```

## Benchmarks
```text
BenchmarkNewDefault-48     	 3656020	       332 ns/op	     392 B/op	       4 allocs/op
BenchmarkLogEmpty-48       	19264894	        60.5 ns/op	      40 B/op	       2 allocs/op
BenchmarkDisabled-48       	1000000000	         0.164 ns/op	       0 B/op	       0 allocs/op
BenchmarkInfo-48           	17027442	        64.3 ns/op	      40 B/op	       2 allocs/op
BenchmarkWithFields-48     	18255840	        65.5 ns/op	      40 B/op	       2 allocs/op
BenchmarkLogFields-48      	12839894	        88.3 ns/op	      40 B/op	       2 allocs/op
BenchmarkLog10Fields-48    	 8441210	       128 ns/op	      75 B/op	       4 allocs/op
BenchmarkLog10String-48    	14184310	        90.9 ns/op	      40 B/op	       2 allocs/op
BenchmarkLog10Int64-48     	10568299	        99.8 ns/op	      40 B/op	       2 allocs/op
``` 
 
## References

Inspired by following projects:

- [uber/zap/buffer](https://go.uber.org/zap/buffer)
- [rs/zerolog](https://github.com/rs/zerolog)
- [qingstor/log](https://github.com/qingstor/log)
