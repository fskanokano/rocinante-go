package bind

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
)

type JSON struct {
	Data     []byte
	Validate *validator.Validate
}

func (j JSON) Bind(s interface{}) error {
	if err := json.Unmarshal(j.Data, s); err != nil {
		return err
	}

	return j.Validate.Struct(s)
}
