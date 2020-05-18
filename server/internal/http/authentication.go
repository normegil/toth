package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/google/uuid"
	"github.com/normegil/toth/server/internal"
	httperror "github.com/normegil/toth/server/internal/http/error"
	"github.com/normegil/toth/server/internal/security"
	"net/http"
	"time"
)

type UserDAO interface {
	security.UserDAO
	Load(id uuid.UUID) (*internal.User, error)
}

type AuthenticationDependencies struct {
	ErrHandler     httperror.HTTPErrorHandler
	UserDAO        UserDAO
	SessionManager *scs.SessionManager
}

func NewAuthenticationMiddleware(handler http.Handler, dependencies AuthenticationDependencies) http.Handler {
	handler = SessionHandler{
		Handler:              handler,
		SessionManager:       dependencies.SessionManager,
		ErrHandler:           dependencies.ErrHandler,
		UserDAO:              dependencies.UserDAO,
		RequestAuthenticator: newRequestAuthenticator(dependencies.UserDAO, dependencies.SessionManager),
	}
	return anonymousUserSetter{Handler: handler}
}

func newRequestAuthenticator(userDAO security.UserDAO, sessionManager *scs.SessionManager) requestAuthenticator {
	authenticator := security.Authenticator{DAO: userDAO}
	updater := AuthenticatedUserSessionUpdater{
		SessionManager: sessionManager,
	}
	requestAuthenticator := requestAuthenticator{
		UserValidator:   authenticator,
		OnAuthenticated: updater.RenewSessionOnAuthenticatedUser,
	}
	return requestAuthenticator
}

const KeyUser string = "authenticated-user"

type anonymousUserSetter struct {
	Handler http.Handler
}

func (a anonymousUserSetter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(context.WithValue(r.Context(), KeyUser, internal.UserAnonymous()))
	a.Handler.ServeHTTP(w, r)
}

type UserValidator interface {
	Authenticate(mail internal.Mail, password string) (*internal.User, error)
}

type requestAuthenticator struct {
	UserValidator   UserValidator
	OnAuthenticated func(r *http.Request, userID uuid.UUID) error
}

func (a requestAuthenticator) Authenticate(r *http.Request) error {
	mailStr, password, ok := r.BasicAuth()
	if !ok {
		return nil
	}
	mail, err := internal.NewMail(mailStr)
	if err != nil {
		return err
	}
	user, err := a.UserValidator.Authenticate(mail, password)
	if nil != err && !security.IsInvalidAuthentication(err) {
		return fmt.Errorf("error during authentication: %w", err)
	}
	if nil != user {
		if nil != a.OnAuthenticated {
			if err := a.OnAuthenticated(r, user.ID); nil != err {
				return fmt.Errorf("authenticater user event error: %w", err)
			}
		}
	} else {
		return httperror.HTTPError{
			Code:   40100,
			Status: http.StatusUnauthorized,
			Err:    errors.New("authentication failed: wrong user and/or password"),
		}
	}
	return nil
}

const keySessionUser = "user"

type SessionHandler struct {
	SessionManager       *scs.SessionManager
	RequestAuthenticator requestAuthenticator
	ErrHandler           httperror.HTTPErrorHandler
	UserDAO              UserDAO
	Handler              http.Handler
}

func (s SessionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var token string
	cookie, err := r.Cookie(s.SessionManager.Cookie.Name)
	if err == nil {
		token = cookie.Value
	}

	ctx, err := s.SessionManager.Load(r.Context(), token)
	if err != nil {
		s.ErrHandler.Handle(w, fmt.Errorf("could not load session: %w", err))
		return
	}
	r = r.WithContext(ctx)

	if err := s.handleAuthenticationAction(r); nil != err {
		s.ErrHandler.Handle(w, err)
		return
	}

	if err := s.RequestAuthenticator.Authenticate(r); nil != err {
		s.ErrHandler.Handle(w, err)
		return
	}

	userID := s.SessionManager.Get(ctx, keySessionUser)
	if nil != userID {
		userIDStr := userID.(string)
		if "" != userIDStr && internal.UserAnonymous().ID.String() != userIDStr {
			parsedID, err := uuid.Parse(userIDStr)
			if err != nil {
				s.ErrHandler.Handle(w, fmt.Errorf("cannot parse id %s: %w", userIDStr, err))
				return
			}
			user, err := s.UserDAO.Load(parsedID)
			if err != nil {
				s.ErrHandler.Handle(w, fmt.Errorf("could not load user '%s': %w", userIDStr, err))
				return
			}
			ctx = context.WithValue(ctx, KeyUser, *user)
		}
	}
	sr := r.WithContext(ctx)

	switch s.SessionManager.Status(ctx) {
	case scs.Unmodified:
		fallthrough
	case scs.Modified:
		token, expiry, err := s.SessionManager.Commit(ctx)
		if err != nil {
			s.ErrHandler.Handle(w, fmt.Errorf("could not commit session: %w", err))
			return
		}
		s.writeSession(w, token, expiry)
	case scs.Destroyed:
		s.writeSession(w, "", time.Time{})
	}

	s.Handler.ServeHTTP(w, sr)
}

func (s SessionHandler) handleAuthenticationAction(r *http.Request) error {
	authenticationAction := r.Header.Get("X-Authentication-Action")
	if authenticationAction != "" {
		userSessionUpdater := AuthenticatedUserSessionUpdater{SessionManager: s.SessionManager}
		switch authenticationAction {
		case "sign-out":
			err := userSessionUpdater.SignOut(r)
			if err != nil {
				return fmt.Errorf("couldn't sign out: %w", err)
			}
		default:
			return fmt.Errorf("unrecognized authentication action: '%s'", authenticationAction)
		}
	}
	return nil
}

func (s SessionHandler) writeSession(w http.ResponseWriter, token string, expiry time.Time) {
	cookie := &http.Cookie{
		Name:     s.SessionManager.Cookie.Name,
		Value:    token,
		Path:     s.SessionManager.Cookie.Path,
		Domain:   s.SessionManager.Cookie.Domain,
		Secure:   s.SessionManager.Cookie.Secure,
		HttpOnly: s.SessionManager.Cookie.HttpOnly,
		SameSite: s.SessionManager.Cookie.SameSite,
	}

	if expiry.IsZero() {
		cookie.Expires = time.Unix(1, 0)
		cookie.MaxAge = -1
	} else if s.SessionManager.Cookie.Persist {
		cookie.Expires = time.Unix(expiry.Unix()+1, 0)        // Round up to the nearest second.
		cookie.MaxAge = int(time.Until(expiry).Seconds() + 1) // Round up to the nearest second.
	}

	w.Header().Add("Set-Cookie", cookie.String())
	addHeaderIfMissing(w, "Cache-Control", `no-cache="Set-Cookie"`)
	addHeaderIfMissing(w, "Vary", "Cookie")
}

func addHeaderIfMissing(w http.ResponseWriter, key, value string) {
	for _, h := range w.Header()[key] {
		if h == value {
			return
		}
	}
	w.Header().Add(key, value)
}

type AuthenticatedUserSessionUpdater struct {
	SessionManager *scs.SessionManager
}

func (a AuthenticatedUserSessionUpdater) RenewSessionOnAuthenticatedUser(r *http.Request, userID uuid.UUID) error {
	if err := a.SessionManager.RenewToken(r.Context()); nil != err {
		return fmt.Errorf("could not renew session token: %w", err)
	}
	a.SessionManager.Put(r.Context(), keySessionUser, userID.String())
	return nil
}

func (a AuthenticatedUserSessionUpdater) SignOut(r *http.Request) error {
	if err := a.SessionManager.RenewToken(r.Context()); nil != err {
		return fmt.Errorf("could not renew session token: %w", err)
	}
	a.SessionManager.Put(r.Context(), keySessionUser, internal.UserAnonymous().ID.String())
	return nil
}
