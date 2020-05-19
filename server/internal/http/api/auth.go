package api

import (
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	internalhttp "github.com/normegil/toth/server/internal/http"
	httperror "github.com/normegil/toth/server/internal/http/error"
	"net/http"
)

type Auth struct {
	ErrHandler     httperror.HTTPErrorHandler
	SessionUpdater internalhttp.AuthenticatedUserSessionUpdater
}

func NewAuth(errHandler httperror.HTTPErrorHandler, sessionManager *scs.SessionManager) *Auth {
	return &Auth{
		ErrHandler: errHandler,
		SessionUpdater: internalhttp.AuthenticatedUserSessionUpdater{
			SessionManager: sessionManager,
		},
	}
}

func (a Auth) Handler() http.Handler {
	r := chi.NewRouter()
	r.Get("/log-in", NeutralResponse{}.Handle)
	r.Get("/log-out", a.LogOut)
	return r
}

func (a Auth) LogOut(w http.ResponseWriter, r *http.Request) {
	err := a.SessionUpdater.SignOut(r)
	if err != nil {
		a.ErrHandler.Handle(w, err)
		return
	}
}
