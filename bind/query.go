package bind

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/url"
)

type Query struct {
	Data     url.Values
	Validate *validator.Validate
}

func (q Query) getQueryJSONByte() ([]byte, error) {
	m := make(map[string]string)

	for key, values := range q.Data {
		if len(values) >= 1 {
			m[key] = values[0]
		}
	}

	return json.Marshal(m)
}

func (q Query) Bind(s interface{}) error {
	//queryJSONByte, err := q.getQueryJSONByte()
	//if err != nil {
	//	return err
	//}
	//
	//if err = json.Unmarshal(queryJSONByte, s); err != nil {
	//	return err
	//}
	if err := unmarshalValues(s, q.Data, "query"); err != nil {
		return err
	}
	return q.Validate.Struct(s)
}
