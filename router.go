package rocinante

import (
	"fmt"
	"net/http"
)

type router struct {
	*trie
	cache    *LRUCache
	handlers map[string][]Handler
}

func newRouter() *router {
	return &router{
		trie:     newTrie(),
		cache:    defaultCache(),
		handlers: make(map[string][]Handler),
	}
}

func (r *router) handle(c *Context) {
	cacheKey := c.Request.URL.Path + "-" + c.Request.Method
	value, ok := r.cache.Get(cacheKey)
	if ok {
		handlers, params := value.(cacheValue).getValue()
		c.params = params
		c.handlers = handlers
		c.Next()
		return
	}

	handlers, params, ok := r.getRoute(c.Request.URL.Path, c.Request.Method)
	if !ok {
		c.handlers = []Handler{
			DefaultLogger(), NotFoundHandler,
		}
		c.Next()
		return
	}

	c.params = params
	c.handlers = handlers
	c.Next()

	if c.StatusCode != http.StatusNotFound {
		r.cache.Set(cacheKey, newCacheValue(params, handlers))
	}
}

func (r *router) addRoute(pattern string, method string, handlers ...Handler) {
	key := fmt.Sprintf("%s-%s", pattern, method)
	if _, exists := r.handlers[key]; exists {
		panic(fmt.Sprintf(`pattern "%s" is already registered in method "%s"`, pattern, method))
	}
	r.insert(pattern, method)
	r.handlers[key] = handlers
}

func (r *router) getRoute(pattern string, method string) ([]Handler, Params, bool) {
	matchedPattern, params, ok := r.search(pattern, method)
	if ok {
		return r.handlers[fmt.Sprintf("%s-%s", matchedPattern, method)], params, ok
	} else {
		return nil, nil, false
	}
}

type cacheValue struct {
	params   Params
	handlers []Handler
}

func (v cacheValue) getValue() ([]Handler, Params) {
	return v.handlers, v.params
}

func newCacheValue(params Params, handlers []Handler) cacheValue {
	return cacheValue{
		params:   params,
		handlers: handlers,
	}
}
