package router

import (
	"github.com/go-chi/chi"
	"net/http"
)

type Controller interface {
	Path() string
	Handler() http.Handler
}

var static http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})

type Dependencies struct {
	APIControllers []Controller
}

func New(d Dependencies) http.Handler {
	r := chi.NewRouter()
	r.Handle("/", static)
	r.Handle("/api", toHandler(d.APIControllers))
	return r
}

func toHandler(ctrls []Controller) http.Handler {
	r := chi.NewRouter()
	for _, ctrl := range ctrls {
		r.Handle(ctrl.Path(), ctrl.Handler())
	}
	return r
}
