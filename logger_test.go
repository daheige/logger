package logger

import (
	"context"
	"log"
	"runtime/debug"
	"testing"
	"time"

	"go.uber.org/zap"
)

// TestLogger test logger.
func TestLogger(t *testing.T) {
	// 对于option 下面的可以根据实际情况使用
	var logger = New(
		WithLogDir("./logs"),
		WithLogFilename("zap.log"),
		WithStdout(true), // 一般生产环境，建议不输出到stdout
		WithJsonFormat(true),
		WithAddCaller(true),
		WithCallerSkip(1), // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		WithEnableColor(false),
		WithLogLevel(zap.DebugLevel), // 设置日志打印最低级别,如果不设置默认为info级别
		WithMaxAge(3),
		WithMaxSize(20),
		WithCompress(false),
		WithHostname("myapp.com"),
		WithEnableCatchStack(true), // 当使用Panic方法时候是否记录stack信息
	)

	// reqId := RndUUID()
	reqId := RndUUIDMd5()
	ctx := context.Background()
	ctx = context.WithValue(ctx, XRequestID, reqId)
	logger.Info(ctx, "hello", map[string]interface{}{
		"a": 1,
		"b": 12,
	})

	logger.Error(ctx, "exec error", zap.Any("details", map[string]interface{}{
		"name": "zap",
		"age":  30,
	}))

	logger.Debug(ctx, "test abc", nil)

	logger.Warn(ctx, "run warning", "key", 12)
	logger.DPanic(ctx, "exec panic but not exit", "stack", string(debug.Stack()))

	logger.Info(ctx, "abc")

	go func() {
		defer logger.CatchPanic(ctx, "exec panic", "key", 123)

		x := 1
		log.Println("x = ", x)
		// panic(1111)
		logger.Panic(ctx, "current goroutine exit")

	}()

	time.Sleep(3 * time.Second)
	log.Println("exit...")
}

// TestNewLogSugar test log sugar.
func TestNewLogSugar(t *testing.T) {
	// 测试log sugar方法
	logSugar := NewLogSugar(WithLogDir("./logs"),
		WithLogFilename("zap-sugar.log"),
		WithStdout(true), // 一般生产环境，建议不输出到stdout
		WithJsonFormat(true),
		WithAddCaller(true),
		WithCallerSkip(1), // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		WithEnableColor(false),
		WithLogLevel(zap.DebugLevel), // 设置日志打印最低级别,如果不设置默认为info级别
		WithMaxAge(3),
		WithMaxSize(20),
		WithCompress(false),
		WithHostname("myapp.com"),
		WithEnableCatchStack(true), // 当使用Panic方法时候是否记录stack信息)
	)

	logSugar.Info("abc", 123, "info", "sugar hello")
	logSugar.Error("a", 234, "x", "sugar hello world")
}

func BenchmarkNew(b *testing.B) {
	// 对于option 下面的可以根据实际情况使用
	var logger = New(
		WithLogDir("./logs"),
		WithLogFilename("zap-bench.log"),
		WithStdout(true), // 一般生产环境，建议不输出到stdout
		WithJsonFormat(true),
		WithAddCaller(true),
		WithCallerSkip(1), // 如果基于这个Logger包，再包装一次，这个skip = 2,以此类推
		WithEnableColor(false),
		WithLogLevel(zap.DebugLevel), // 设置日志打印最低级别,如果不设置默认为info级别
		WithMaxAge(3),
		WithMaxSize(20),
		WithCompress(false),
		// WithHostname("myapp.com"),
		WithEnableCatchStack(true), // 当使用Panic方法时候是否记录stack信息
	)

	// reqId := RndUUID()
	reqId := RndUUIDMd5()
	ctx := context.Background()
	ctx = context.WithValue(ctx, XRequestID, reqId)
	logger.Info(ctx, "exec begin")
	start := time.Now()
	for i := 0; i < b.N; i++ {
		logger.Info(ctx, "hello", "index", i)
		logger.Error(ctx, "exec error", "abc", 1, "e", "zap is fast")
		logger.Info(ctx, "exec map", map[string]interface{}{
			"a": 1,
			"b": 123.23,
			"c": "hello,go",
			"e": []string{"f", "g", "higk"},
			"f": []int{1, 2, 3, i},
		})
	}

	logger.Info(ctx, "exec end", "cost_time", time.Since(start).Seconds())
}

/**
BenchmarkNew test
{"level":"info","time_local":"2020-09-20T17:21:35.883+0800",
"caller_line":"/Users/heige/web/go/logger/logger_test.go:116",
"msg":"exec map","a":1,"b":123.23,"c":"hello,go",
"e":["f","g","higk"],"f":[1,2,3,18444],"local_time":"2020-09-20 17:21:35.883",
"hostname":"daheige","x-request-id":"4bbc721e9da802cfee4cdfb3689220e0"}
{"level":"info","time_local":"2020-09-20T17:21:35.884+0800",
"caller_line":"/Users/heige/web/go/logger/logger_test.go:125",
"msg":"exec end","cost_time":1.338923502,"local_time":"2020-09-20 17:21:35.884",
"hostname":"daheige","x-request-id":"4bbc721e9da802cfee4cdfb3689220e0"}
BenchmarkNew-12    	   18445	     72602 ns/op
PASS
*/
