package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type customHttpResponse struct {
	StatusCode int         `json:"statusCode"`
	Message    interface{} `json:"message"`
}

func responseBuilder(w http.ResponseWriter, code int, msg ...interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)

	if msg == nil || code == http.StatusInternalServerError {
		return
	}

	resp := customHttpResponse{
		StatusCode: code,
	}

	switch v := msg[0].(type) {
	case error:
		resp.Message = v.Error()
	case []byte:
		resp.Message = string(v)
	default:
		resp.Message = v
	}

	b, _ := json.Marshal(resp)
	w.Write(b)
}

// readData reads incoming body and checks for specified fields in it
func readData(req *http.Request, fields ...string) (map[string]interface{}, error) {
	if req.Body == nil {
		return nil, errors.New("no data in body")
	}
	var input map[string]interface{}
	if err := json.NewDecoder(req.Body).Decode(&input); err != nil {
		if err == io.EOF {
			return nil, err
		}
		return nil, err
	}
	return input, checkData(input, fields...)
}

func checkData(m map[string]interface{}, fields ...string) error {
	if fields != nil {
		for _, v := range fields {
			if _, ok := m[v]; !ok {
				return fmt.Errorf("field %s must be exists", v)
			}
		}
	}
	return nil
}
