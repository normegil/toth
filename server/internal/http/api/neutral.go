package api

import "net/http"

type NeutralResponse struct {
}

func (_ NeutralResponse) Handle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNoContent)
}
