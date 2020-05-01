package error

import (
	"encoding/json"
	"errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type HTTPErrorHandler struct {
	LogUserError bool
}

func (h HTTPErrorHandler) Handle(w http.ResponseWriter, err error) {
	httpError := &HTTPError{}
	if !errors.As(err, httpError) {
		httpError = &HTTPError{
			Code:   50000,
			Status: http.StatusInternalServerError,
			Err:    err,
		}
	}
	if httpError.Status%400 < 100 && h.LogUserError {
		log.Error().Err(err).Msg("user error")
	} else if httpError.Status%400 >= 100 {
		log.Error().Err(err).Msg("rest error")
	}
	resp := httpError.toResponse()
	resp.Error = err.Error()
	bytes, marshalErr := json.Marshal(resp)
	if marshalErr != nil {
		log.Error().Err(marshalErr).Interface("HTTPError", resp).Msg("Could not marshal HTTPError")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.Status)
	if _, writeErr := w.Write(bytes); nil != writeErr {
		log.Error().Err(writeErr).Interface("HTTPError", resp).Msg("Could not write response with HTTPError")
	}
}
