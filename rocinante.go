package rocinante

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
	"sync"
)

type Rocinante struct {
	*RouterGroup
	groups      []*RouterGroup
	router      *router
	contextPool sync.Pool
	Validate    *validator.Validate
}

func New() *Rocinante {
	app := &Rocinante{
		router:      newRouter(),
		contextPool: newContextPool(),
		Validate:    validator.New(),
	}
	app.RouterGroup = newRouterGroup("", app)
	app.groups = []*RouterGroup{app.RouterGroup}
	initLogger()
	return app
}

func Default() *Rocinante {
	app := New()
	app.Use(DefaultLogger(), Recovery())
	return app
}

func newContextPool() sync.Pool {
	return sync.Pool{New: func() interface{} {
		return newContext()
	}}
}

func (r *Rocinante) Run(addr ...string) error {
	resolvedAddr := resolveAddr(addr...)
	r.printStartServerLog(resolvedAddr, "HTTP")
	return http.ListenAndServe(resolvedAddr, r)
}

func (r *Rocinante) RunTLS(certFile string, keyFile string, addr ...string) error {
	resolvedAddr := resolveAddr(addr...)
	r.printStartServerLog(resolvedAddr, "HTTPS")
	return http.ListenAndServeTLS(resolvedAddr, certFile, keyFile, r)
}

func (r *Rocinante) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := r.contextPool.Get().(*Context)
	c.reset(w, req, r)

	r.router.handle(c)

	r.contextPool.Put(c)
}

func (r *Rocinante) SetLogger(l Logger) {
	logger = l
}

func (r *Rocinante) printStartServerLog(addr string, protocol string) {
	logger.Debug("Starting Server ...\n")
	for key, handlers := range r.router.handlers {
		sp := strings.Split(key, "-")
		method := sp[1]
		if method == "OPTIONS" {
			continue
		}
		path := sp[0]
		lh := len(handlers)
		logger.Debug(fmt.Sprintf("%s %s --> %s (%d handlers)", method, path, getFunctionName(handlers[lh-1]), lh))
	}
	logger.Debug(fmt.Sprintf("Listening and serving %s on %s\n", protocol, addr))
}
