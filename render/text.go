package render

import (
	"net/http"
	"unsafe"
)

type String struct {
	Data string
}

func (s String) Render(w http.ResponseWriter, status int) error {
	s.WriteContentType(w)

	writeStatus(w, status)
	_, err := w.Write(StringToBytes(s.Data))
	return err
}

func (s String) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, TextContentType)
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

const TextContentType = "text/plain; charset=utf-8"
