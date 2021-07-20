package render

import (
	"encoding/json"
	"net/http"
)

type JSON struct {
	Data interface{}
}

func (j JSON) Render(w http.ResponseWriter, status int) error {
	j.WriteContentType(w)

	jsonBytes, err := json.Marshal(j.Data)
	if err != nil {
		return err
	}

	writeStatus(w, status)
	_, err = w.Write(jsonBytes)
	return err
}

func (j JSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, JSONContentType)
}

const JSONContentType = "application/json; charset=utf-8"
