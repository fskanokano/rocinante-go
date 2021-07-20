package render

import (
	"net/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, status int) error
	WriteContentType(w http.ResponseWriter)
}

func writeContentType(w http.ResponseWriter, contentType string) {
	w.Header().Set("Content-Type", contentType)
}

func writeStatus(w http.ResponseWriter, status int) {
	w.WriteHeader(status)
}
