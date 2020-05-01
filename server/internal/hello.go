package internal

import (
	httperror "github.com/normegil/toth/server/internal/http/error"
	"net/http"
)

type Hello struct {
	ErrHandler httperror.HTTPErrorHandler
}

func (h Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello")); nil != err {
		h.ErrHandler.Handle(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
