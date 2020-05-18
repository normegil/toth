package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/normegil/toth/server/internal"
	internalhttp "github.com/normegil/toth/server/internal/http"
	httperror "github.com/normegil/toth/server/internal/http/error"
	"net/http"
)

type userRTO struct {
	ID   uuid.UUID
	Name string
	Mail string
}

func fromUser(u internal.User) *userRTO {
	return &userRTO{
		ID:   u.ID,
		Name: u.Name,
		Mail: u.Mail.String(),
	}
}

type Users struct {
	ErrHandler httperror.HTTPErrorHandler
}

func (a Users) Path() string {
	return "/users"
}

func (a Users) Handler() http.Handler {
	r := chi.NewRouter()
	r.Get("/current", http.HandlerFunc(a.currentUser))
	return r
}

func (a Users) currentUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(internalhttp.KeyUser)
	rto := fromUser(user.(internal.User))

	bytes, err := json.Marshal(rto)
	if err != nil {
		a.ErrHandler.Handle(w, fmt.Errorf("marshalling %+v: %w", rto, err))
		return
	}
	if _, err := w.Write(bytes); nil != err {
		a.ErrHandler.Handle(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}
