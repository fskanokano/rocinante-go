package rocinante

import (
	"fmt"
	"net/http"
	"path"
	"strings"
)

type RouterGroup struct {
	prefix      string
	middlewares []Handler
	parent      *RouterGroup
	app         *Rocinante
}

func newRouterGroup(prefix string, app *Rocinante) *RouterGroup {
	return &RouterGroup{
		prefix: prefix,
		app:    app,
	}
}

func (r *RouterGroup) GET(relativePath string, handlers ...Handler) {
	r.addRoute(relativePath, "GET", false, handlers...)
}

func (r *RouterGroup) POST(relativePath string, handlers ...Handler) {
	r.addRoute(relativePath, "POST", false, handlers...)
}

func (r *RouterGroup) PUT(relativePath string, handlers ...Handler) {
	r.addRoute(relativePath, "PUT", false, handlers...)
}

func (r *RouterGroup) DELETE(relativePath string, handlers ...Handler) {
	r.addRoute(relativePath, "DELETE", false, handlers...)
}

func (r *RouterGroup) PATCH(relativePath string, handlers ...Handler) {
	r.addRoute(relativePath, "PATCH", false, handlers...)
}

func (r *RouterGroup) OPTIONS(relativePath string, handlers ...Handler) {
	r.addRoute(relativePath, "OPTIONS", false, handlers...)
}

func (r *RouterGroup) HEAD(relativePath string, handlers ...Handler) {
	r.addRoute(relativePath, "HEAD", false, handlers...)
}

func (r *RouterGroup) Route(relativePath string, controller ControllerInterface) {
	dispatchHandler(relativePath, controller, r)
}

func (r *RouterGroup) Group(prefix string) *RouterGroup {
	if !strings.HasPrefix(prefix, "/") {
		panic("invalid prefix")
	}

	prefix = r.prefix + prefix
	newGroup := newRouterGroup(prefix, r.app)
	newGroup.parent = r
	newGroup.middlewares = r.middlewares

	r.app.groups = append(r.app.groups, newGroup)
	return newGroup
}

func (r *RouterGroup) Use(middlewares ...Handler) {
	r.middlewares = append(r.middlewares, middlewares...)
}

func (r *RouterGroup) Static(relativePath string, root string) {
	handler := r.staticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	r.GET(urlPattern, handler)
}

func (r *RouterGroup) WebSocket(relativePath string, handler WebsocketHandler, middlewares ...Handler) {
	var websocketHandler = func(c *Context) {
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			panic(err)
		}
		handler(conn)
	}
	handlers := append(middlewares, websocketHandler)
	r.GET(relativePath, handlers...)
}

func (r *RouterGroup) staticHandler(relativePath string, fs http.FileSystem) Handler {
	absolutePath := path.Join(r.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			NotFoundHandler(c)
			return
		}

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (r *RouterGroup) addRoute(relativePath string, method string, isMVC bool, handlers ...Handler) {
	if len(handlers) < 1 {
		panic("invalid handlers")
	}

	fullPath := r.prefix + relativePath
	r.app.router.addRoute(fullPath, method, r.finalHandlers(handlers...)...)

	if !isMVC {
		if method != "OPTIONS" && method != "GET" {
			optionsHandlerKey := fmt.Sprintf("%s-%s", fullPath, "OPTIONS")
			if _, exists := r.app.router.handlers[optionsHandlerKey]; !exists {
				r.app.router.addRoute(fullPath, "OPTIONS", r.finalHandlers(OptionsHandler)...)
			}
		}
	}
}

func (r *RouterGroup) finalHandlers(handlers ...Handler) []Handler {
	final := make([]Handler, 0)
	for _, middleware := range r.middlewares {
		final = append(final, middleware)
	}
	for _, handler := range handlers {
		final = append(final, handler)
	}

	if len(final) >= int(abortIndex) {
		panic("too many handlers")
	}

	return final
}
