package rocinante

import (
	"bytes"
	"github.com/fskanokano/rocinante-go/bind"
	"github.com/fskanokano/rocinante-go/render"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

type Context struct {
	Writer     http.ResponseWriter
	Request    *http.Request
	params     Params
	session    Session
	StatusCode int
	Method     string
	handlers   HandlersChain
	index      int8
	app        *Rocinante
}

func newContext() *Context {
	return &Context{StatusCode: 200}
}

func (c *Context) reset(w http.ResponseWriter, req *http.Request, app *Rocinante) {
	c.Writer = w
	c.Request = req
	c.session = make(Session)
	c.Method = req.Method
	c.index = -1
	c.StatusCode = 200
	if c.app == nil {
		c.app = app
	}
}

func (c *Context) JSON(v interface{}, status ...int) {
	c.Render(render.JSON{Data: v}, status...)
}

func (c *Context) String(s string, status ...int) {
	c.Render(render.String{Data: s}, status...)
}

func (c *Context) HTML(name string, data interface{}, status ...int) {
	c.Render(render.HTML{
		Name: name,
		Data: data,
	}, status...)
}

func (c *Context) File(filepath string) {
	http.ServeFile(c.Writer, c.Request, filepath)
}

func (c *Context) Render(r render.Renderer, status ...int) {
	resolvedStatus := resolveStatus(status...)
	c.StatusCode = resolvedStatus
	if err := r.Render(c.Writer, resolvedStatus); err != nil {
		panic(err)
	}
}

func (c *Context) BindJSON(s interface{}) error {
	raw, err := c.GetRawData()
	if err != nil {
		return err
	}

	return c.Bind(bind.JSON{
		Data:     raw,
		Validate: c.app.Validate,
	}, s)
}

func (c *Context) BindForm(s interface{}) error {
	if c.Request.Form == nil {
		c.initForm()
	}

	return c.Bind(bind.Form{
		Data:     c.Request.Form,
		Validate: c.app.Validate,
	}, s)
}

func (c *Context) BindQuery(s interface{}) error {
	return c.Bind(bind.Query{
		Data:     c.Request.URL.Query(),
		Validate: c.app.Validate,
	}, s)
}

func (c *Context) Bind(binder bind.Binder, s interface{}) error {
	return binder.Bind(s)
}

func (c *Context) SetResponseHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) GetResponseHeader(key string) string {
	return c.Writer.Header().Get(key)
}

func (c *Context) AddResponseHeader(key string, value string) {
	c.Writer.Header().Add(key, value)
}

func (c *Context) DelResponseHeader(key string) {
	c.Writer.Header().Del(key)
}

func (c *Context) SetRequestHeader(key string, value string) {
	c.Request.Header.Set(key, value)
}

func (c *Context) GetRequestHeader(key string) string {
	return c.Request.Header.Get(key)
}

func (c *Context) AddRequestHeader(key string, value string) {
	c.Request.Header.Add(key, value)
}

func (c *Context) DelRequestHeader(key string) {
	c.Request.Header.Del(key)
}

func (c *Context) SetForm(key string, value string) {
	if c.Request.Form == nil {
		c.initForm()
	}
	c.Request.Form.Set(key, value)
}

func (c *Context) GetForm(key string) string {
	if c.Request.Form == nil {
		c.initForm()
	}
	return c.Request.Form.Get(key)
}

func (c *Context) AddForm(key string, value string) {
	if c.Request.Form == nil {
		c.initForm()
	}
	c.Request.Form.Add(key, value)
}

func (c *Context) DelForm(key string) {
	if c.Request.Form == nil {
		c.initForm()
	}
	c.Request.Form.Del(key)
}

func (c *Context) initForm() {
	c.Request.ParseMultipartForm(defaultMaxMemory)
}

func (c *Context) SetQuery(key string, value string) {
	c.Request.URL.Query().Set(key, value)
}

func (c *Context) GetQuery(key string) string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context) AddQuery(key string, value string) {
	c.Request.URL.Query().Add(key, value)
}

func (c *Context) DelQuery(key string) {
	c.Request.URL.Query().Del(key)
}

func (c *Context) FormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return c.Request.FormFile(key)
}

func (c *Context) FormValue(key string) string {
	return c.Request.FormValue(key)
}

func (c *Context) DefaultForm(key string, defaultValue string) string {
	query := c.GetForm(key)
	if query != "" {
		return query
	} else {
		return defaultValue
	}
}

func (c *Context) DefaultQuery(key string, defaultValue string) string {
	query := c.GetQuery(key)
	if query != "" {
		return query
	} else {
		return defaultValue
	}
}

func (c *Context) SetStatus(status int) {
	c.StatusCode = status
	c.Writer.WriteHeader(status)
}

func (c *Context) Param(key string) string {
	return c.params[key]
}

func (c *Context) SetSession(key string, value interface{}) {
	c.session[key] = value
}

func (c *Context) GetSession(key string) interface{} {
	return c.session[key]
}

func (c *Context) ReadRawData() ([]byte, error) {
	return ioutil.ReadAll(c.Request.Body)
}

func (c *Context) GetRawData() ([]byte, error) {
	rawData, err := c.ReadRawData()
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(rawData))
	return rawData, err
}

func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *Context) Abort() {
	c.index = abortIndex
}

func (c *Context) AbortWithJSON(v interface{}, status ...int) {
	c.JSON(v, status...)
	c.Abort()
}

func (c *Context) AbortWithString(s string, status ...int) {
	c.String(s, status...)
	c.Abort()
}

func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.Request.Header.Get("Connection")), "upgrade") &&
		strings.EqualFold(c.Request.Header.Get("Upgrade"), "websocket") {
		return true
	}
	return false
}

const abortIndex int8 = math.MaxInt8 / 2

type Params map[string]string

type Session map[string]interface{}

type Map map[string]interface{}

type Handler func(*Context)

type HandlersChain []Handler
