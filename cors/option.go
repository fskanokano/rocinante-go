package cors

import (
	"strconv"
	"strings"
)

type Option struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	AllowCredentials bool
	ExposeHeaders    []string
	MaxAge           int
}

func DefaultOption() Option {
	return Option{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		ExposeHeaders:    []string{"*"},
		MaxAge:           600,
	}
}

type optionValue struct {
	allowOrigins     string
	allowMethods     string
	allowHeaders     string
	allowCredentials bool
	exposeHeaders    string
	maxAge           string
}

func (o *optionValue) parseOption(option Option) {
	if option.AllowMethods == nil || option.AllowMethods[0] == "*" {
		o.allowMethods = AllMethods
	} else {
		o.allowMethods = o.combineOptionSlice(option.AllowMethods)
	}

	if option.AllowHeaders == nil || option.AllowHeaders[0] == "*" {
		o.allowHeaders = SafeHeaders
	} else {
		o.allowHeaders = o.combineOptionSlice(option.AllowHeaders)
	}

	if option.ExposeHeaders == nil || option.ExposeHeaders[0] == "*" {
		o.exposeHeaders = SafeHeaders
	} else {
		o.exposeHeaders = o.combineOptionSlice(option.ExposeHeaders)
	}

	o.allowCredentials = option.AllowCredentials

	if option.MaxAge == 0 {
		o.maxAge = "600"
	} else {
		o.maxAge = strconv.Itoa(option.MaxAge)
	}
}

func (o *optionValue) resolveOrigin(origin string, option Option) {
	if option.AllowOrigins == nil || option.AllowOrigins[0] == "*" {
		if o.allowCredentials {
			o.allowOrigins = origin
		} else {
			o.allowOrigins = "*"
		}
	} else {
		o.allowOrigins = origin
	}
}

func (o *optionValue) combineOptionSlice(s []string) string {
	return strings.Join(s, ", ")
}

func (o *optionValue) isAllowedOrigin(origin string, allowOrigins []string) bool {
	for _, allowOrigin := range allowOrigins {
		if origin == allowOrigin {
			return true
		}
	}
	return false
}

const (
	AllMethods  = "DELETE, GET, OPTIONS, PATCH, POST, PUT"
	SafeHeaders = "Accept, Accept-Language, Content-Language, Content-Type"
)
