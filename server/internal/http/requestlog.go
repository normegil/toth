package http

import (
	"github.com/rs/zerolog/log"
	"net/http"
)

type requestLogger struct {
	Handler http.Handler
}

func (l requestLogger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("method", r.Method).Str("url", r.RequestURI).Msg("request")
	l.Handler.ServeHTTP(w, r)
}
