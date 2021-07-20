package bind

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/url"
)

type Form struct {
	Data     url.Values
	Validate *validator.Validate
}

func (f Form) getFormJSONByte() ([]byte, error) {
	m := make(map[string]string)

	for key, values := range f.Data {
		if len(values) >= 1 {
			m[key] = values[0]
		}
	}

	return json.Marshal(m)
}

func (f Form) Bind(s interface{}) error {
	//formJSONByte, err := f.getFormJSONByte()
	//if err != nil {
	//	return err
	//}
	//
	//if err = json.Unmarshal(formJSONByte, s); err != nil {
	//	return err
	//}
	if err := unmarshalValues(s, f.Data, "form"); err != nil {
		return err
	}
	return f.Validate.Struct(s)
}
