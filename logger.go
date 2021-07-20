package rocinante

import (
	"fmt"
	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/fskanokano/rocinante-go/log"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var logger Logger

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warning(args ...interface{})
	Error(args ...interface{})
}

func defaultFormatter(isConsole bool) *nested.Formatter {
	f := &nested.Formatter{
		HideKeys:        true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerFirst:     true,
		ShowFullLevel:   true,
		CustomCallerFormatter: func(frame *runtime.Frame) string {
			funcInfo := runtime.FuncForPC(frame.PC)
			if funcInfo == nil {
				return "error during runtime.FuncForPC"
			}
			fullPath, line := funcInfo.FileLine(frame.PC)
			return fmt.Sprintf(" [%v:%v]", filepath.Base(fullPath), line)
		},
	}
	if isConsole {
		f.NoColors = false
	} else {
		f.NoColors = true
	}
	return f
}

func initLogger() {
	defaultLogger := logrus.New()
	//defaultLogger.SetFormatter(defaultFormatter(true))
	defaultLogger.SetFormatter(log.DefaultFormatter())
	//defaultLogger.SetFormatter(&log.Formatter{ForceColor: false})
	defaultLogger.SetLevel(logrus.DebugLevel)
	logger = defaultLogger
}

func DefaultLogger() Handler {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		costTime := time.Since(start)
		go printDebugLog(costTime, c)
	}
}

func printDebugLog(costTime time.Duration, c *Context) {
	status := c.StatusCode
	host := strings.Split(c.Request.RemoteAddr, ":")[0]
	var method string
	if !c.IsWebsocket() {
		method = c.Method
	} else {
		method = "WebSocket"
	}
	uri := c.Request.RequestURI
	debugLog := fmt.Sprintf(`%d | %s | %s | %s | "%s"`, status, costTime, host, method, uri)
	if status < 400 {
		logger.Info(debugLog)
	} else {
		logger.Error(debugLog)
	}
}
