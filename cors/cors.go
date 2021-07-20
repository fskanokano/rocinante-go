package cors

import (
	"github.com/fskanokano/rocinante-go"
	"net/http"
)

func New(option ...Option) rocinante.Handler {
	opt := resolveOption(option...)

	var optValue optionValue
	optValue.parseOption(opt)

	return func(c *rocinante.Context) {
		origin := c.GetRequestHeader("Origin")
		if origin == "" {
			return
		} else {
			if opt.AllowOrigins != nil || opt.AllowOrigins[0] != "*" {
				if !optValue.isAllowedOrigin(origin, opt.AllowOrigins) {
					c.SetStatus(http.StatusForbidden)
					c.Abort()
					return
				}
			}
		}

		optValue.resolveOrigin(origin, opt)

		if c.Method == "OPTIONS" {
			preflight(c, optValue)
			c.Abort()
			return
		}

		simple(c, optValue)

		c.Next()
	}
}

func preflight(c *rocinante.Context, optValue optionValue) {
	simple(c, optValue)

	c.SetResponseHeader("Access-Control-Allow-Headers", optValue.allowHeaders)
	c.SetResponseHeader("Access-Control-Allow-Methods", optValue.allowMethods)
	c.SetResponseHeader("Access-Control-Max-Age", optValue.maxAge)

	c.SetStatus(http.StatusNoContent)
}

func simple(c *rocinante.Context, optValue optionValue) {
	c.SetResponseHeader("Access-Control-Allow-Origin", optValue.allowOrigins)
	c.SetResponseHeader("Access-Control-Expose-Headers", optValue.exposeHeaders)

	if optValue.allowCredentials {
		c.SetResponseHeader("Access-Control-Allow-Credentials", "true")
	}
}
