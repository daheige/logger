# sentry上报
## install
Install our Go SDK using go get:
```shell
go get github.com/getsentry/sentry-go
```

## Configure SDK
Import and initialize the Sentry SDK early in your application's setup:
```go
package main

import (
  "log"

  "github.com/getsentry/sentry-go"
)

func main() {
  err := sentry.Init(sentry.ClientOptions{
    Dsn: "your sentry dsn",
  })
  if err != nil {
    log.Fatalf("sentry.Init: %s", err)
  }
}
```

## Verify
The quickest way to verify Sentry in your Go program is to capture a message:
```go
package main

import (
  "log"
  "time"

  "github.com/getsentry/sentry-go"
)

func main() {
  err := sentry.Init(sentry.ClientOptions{
    Dsn: "your sentry dsn",
  })
  if err != nil {
    log.Fatalf("sentry.Init: %s", err)
  }
  
  // Flush buffered events before the program terminates.
  defer sentry.Flush(2 * time.Second)

  // 测试上报
  sentry.CaptureMessage("It works!")
}
```

# 日志输出和sentry上报示例
```go
package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	"github.com/daheige/logger"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)

	logger.Default(
		logger.WithJsonFormat(true),        // 默认json格式化输出
		logger.WithCallerSkip(2),           // 如果基于这个Logger包，需要设置适当的skip
		logger.WithLogLevel(zap.InfoLevel), // 设置日志打印最低级别,如果不设置,默认为info级别
		logger.WithStdout(true),            // 日志默认输出到终端

		logger.WithEnableSentry(true), // 开启sentry上报
	)

	logger.Info(context.Background(), "hello world", "plat", "mac")
	logger.Info(context.Background(), "exec begin", "foo", "abc")
}
```

