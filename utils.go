package rocinante

import (
	"net/http"
	"reflect"
	"runtime"
)

func resolveAddr(addr ...string) string {
	switch len(addr) {
	case 0:
		return ":8000"
	default:
		return addr[0]
	}
}

func resolveStatus(status ...int) int {
	switch len(status) {
	case 0:
		return http.StatusOK
	default:
		return status[0]
	}
}

func NotFoundHandler(c *Context) {
	c.String("404 Not Found", http.StatusNotFound)
}

func MethodNotAllowedHandler(c *Context) {
	c.String("405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func InternalServerErrorHandler(c *Context) {
	c.String("500 Internal Server Error", http.StatusInternalServerError)
}

func OptionsHandler(c *Context) {
	c.SetStatus(http.StatusNoContent)
}

func getFunctionName(function interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
}
