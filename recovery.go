package rocinante

import (
	"fmt"
	"runtime"
	"strings"
)

func trace(err error) string {
	message := err.Error()

	var pcs [32]uintptr
	n := runtime.Callers(4, pcs[:]) // skip first 4 caller

	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}

func Recovery() Handler {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				errStack := trace(err.(error))
				logger.Error(errStack)
				InternalServerErrorHandler(c)
			}
		}()

		c.Next()
	}
}
