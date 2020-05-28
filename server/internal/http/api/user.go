package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/normegil/toth/server/internal"
	internalhttp "github.com/normegil/toth/server/internal/http"
	httperror "github.com/normegil/toth/server/internal/http/error"
	"github.com/normegil/toth/server/internal/security"
	"io/ioutil"
	"net/http"
)

type UserDAO interface {
	Load(uuid.UUID) (*internal.User, error)
	UpdateProfile(internal.User) error
	UpdatePassword(userID uuid.UUID, newPasswordHash []byte, newAlgorithmID uuid.UUID) error
}

type passwordRTO struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type userRTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Mail string `json:"mail"`
}

func fromUser(u internal.User) *userRTO {
	return &userRTO{
		ID:   u.ID.String(),
		Name: u.Name,
		Mail: u.Mail.String(),
	}
}

func (u userRTO) toUser() (*internal.User, error) {
	id, err := uuid.Parse(u.ID)
	if err != nil {
		return nil, err
	}
	mail, err := internal.NewMail(u.Mail)
	if err != nil {
		return nil, err
	}
	return &internal.User{
		ID:   id,
		Name: u.Name,
		Mail: mail,
	}, nil
}

type Users struct {
	ErrHandler httperror.HTTPErrorHandler
	UserDAO    UserDAO
}

func (a Users) Handler() http.Handler {
	r := chi.NewRouter()
	r.Get("/current", a.currentUser)
	r.Put("/", a.updateUser)
	r.Put("/{userID}/password", a.updatePassword)
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
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(bytes); nil != err {
		a.ErrHandler.Handle(w, err)
		return
	}
}

func (a Users) updateUser(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.ErrHandler.Handle(w, fmt.Errorf("could not read body: %w", err))
		return
	}

	var userRTO userRTO
	if err := json.Unmarshal(bodyBytes, &userRTO); nil != err {
		a.ErrHandler.Handle(w, fmt.Errorf("could not unmarshal response: %w", err))
		return
	}

	user, err := userRTO.toUser()
	if err != nil {
		a.ErrHandler.Handle(w, fmt.Errorf("parsing to user (%+v): %w", userRTO, err))
		return
	}

	if err = a.UserDAO.UpdateProfile(*user); nil != err {
		a.ErrHandler.Handle(w, fmt.Errorf("updating profile (%+v): %w", user, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (a Users) updatePassword(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		a.ErrHandler.Handle(w, httperror.HTTPError{
			Code:   http.StatusBadRequest,
			Status: 40003,
			Err:    fmt.Errorf("wrong user id '%s': %w", userIDStr, err),
		})
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		a.ErrHandler.Handle(w, fmt.Errorf("reading body: %w", err))
		return
	}

	var rto passwordRTO
	if err := json.Unmarshal(body, &rto); nil != err {
		a.ErrHandler.Handle(w, fmt.Errorf("parsing body: %w", err))
		return
	}

	user, err := a.UserDAO.Load(userID)
	if err != nil {
		a.ErrHandler.Handle(w, fmt.Errorf("loading user (%s): %w", userID.String(), err))
		return
	}

	currentUserAlgorithm := security.AllHashAlgorithms().FindByID(user.HashAlgorithmID)
	if err = currentUserAlgorithm.Validate(user.PasswordHash, rto.CurrentPassword); nil != err {
		a.ErrHandler.Handle(w, httperror.HTTPError{
			Status: http.StatusUnauthorized,
			Code:   40101,
			Err:    fmt.Errorf("invalid password: %w", err),
		})
		return
	}

	algorithm := security.DefaultHashAlgorithm()
	hash, err := algorithm.HashAndSalt(rto.NewPassword)
	if err != nil {
		a.ErrHandler.Handle(w, fmt.Errorf("hash new password: %w", err))
		return
	}

	if err = a.UserDAO.UpdatePassword(userID, hash, algorithm.ID()); nil != err {
		a.ErrHandler.Handle(w, fmt.Errorf("updating password: %w", err))
		return
	}
}
