package http

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"
)

var Static http.Handler = http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusNotFound)
})

func Listen(addr net.TCPAddr, handler http.Handler) error {
	stopHTTPServer := make(chan os.Signal, 1)
	signal.Notify(stopHTTPServer, os.Interrupt)

	closeHttpServer := launch(addr, middlewares(handler))
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := closeHttpServer(ctx); nil != err {
			log.Fatal().Err(err).Msg("closing server failed")
		}
	}()

	<-stopHTTPServer
	return nil
}

func middlewares(handler http.Handler) http.Handler {
	return requestLogger{handler}
}

func launch(addr net.TCPAddr, handler http.Handler) func(ctx context.Context) error {
	httpAddr := addr.IP.String() + ":" + strconv.Itoa(addr.Port)
	srv := http.Server{Addr: httpAddr, Handler: handler}

	go func() {
		log.Info().Str("address", srv.Addr).Msg("HTTP Server listening")
		if err := srv.ListenAndServe(); nil != err {
			if http.ErrServerClosed != err {
				log.Fatal().Err(err).Str("address", srv.Addr).Msg("HTTP Server listening failed")
			}
		}
	}()
	return func(ctx context.Context) error {
		if err := srv.Shutdown(ctx); nil != err {
			return fmt.Errorf("shutdown http server: %w", err)
		}
		return nil
	}
}
