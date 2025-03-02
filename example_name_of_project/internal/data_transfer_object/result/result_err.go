package result

import (
	"ca-generator/example_name_of_project/internal/common/types/error_with_codes"

	"github.com/goccy/go-json"
)

type ResultErr struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

func NewResultErr(err error) *ResultErr {
	var code int

	if errCode, errErr := error_with_codes.ToErrorWithCode(err); errErr != nil {
		code = int(errCode.GetCode())
	}

	return &ResultErr{
		Error: err.Error(),
		Code:  code,
	}
}

func (r *ResultErr) GetJson() ([]byte, error) {
	return json.Marshal(r)
}
