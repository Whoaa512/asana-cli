package output

import (
	"encoding/json"
	"io"

	"github.com/whoaa512/asana-cli/internal/errors"
)

type JSON struct {
	w io.Writer
}

func NewJSON(w io.Writer) *JSON {
	return &JSON{w: w}
}

func (j *JSON) Print(v any) error {
	enc := json.NewEncoder(j.w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Message  string `json:"message"`
	Code     string `json:"code"`
	ExitCode int    `json:"exit_code"`
}

func (j *JSON) PrintError(err error) error {
	cliErr := errors.AsCLIError(err)

	resp := ErrorResponse{
		Error: ErrorDetail{
			Message:  cliErr.Message,
			Code:     cliErr.Code,
			ExitCode: cliErr.ExitCode,
		},
	}

	enc := json.NewEncoder(j.w)
	enc.SetIndent("", "  ")
	return enc.Encode(resp)
}
