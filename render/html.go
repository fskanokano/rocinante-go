package render

import (
	"net/http"
	"text/template"
)

type HTML struct {
	Name string
	Data interface{}
}

func (h HTML) Render(w http.ResponseWriter, status int) error {
	h.WriteContentType(w)
	writeStatus(w, status)
	t, err := template.ParseFiles(h.Name)
	if err != nil {
		return err
	}
	if err = t.Execute(w, h.Data); err != nil {
		return err
	}
	return nil
}

func (h HTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, HTMLContentType)
}

const HTMLContentType = "text/html; charset=utf-8"
