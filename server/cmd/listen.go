package main

import (
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/normegil/godatabaseversioner"
	"github.com/normegil/postgres"
	"github.com/normegil/toth/server/internal/commands"
	internalhttp "github.com/normegil/toth/server/internal/http"
	"github.com/normegil/toth/server/internal/http/api"
	httperror "github.com/normegil/toth/server/internal/http/error"
	internalpostgres "github.com/normegil/toth/server/internal/postgres"
	"github.com/normegil/toth/server/internal/postgres/versions"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net"
)

func listen() (*cobra.Command, error) {
	listenCmd := &cobra.Command{
		Use:   "listen",
		Short: "Launch server",
		Long:  `Launch server`,
		Run:   listenRun,
	}

	listenCmd.Flags().String("listen-address", "0.0.0.0", "Address on which current server will listen")
	if err := viper.BindPFlag("server.address", listenCmd.Flags().Lookup("listen-address")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "server.address", err)
	}

	listenCmd.Flags().Int("listen-port", 8080, "Port on which current server will listen")
	if err := viper.BindPFlag("server.port", listenCmd.Flags().Lookup("listen-port")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "server.port", err)
	}

	listenCmd.Flags().Bool("log-user-error", false, "When logging errors, log errors that were triggered by user (wrong input, wrong call, ...)")
	if err := viper.BindPFlag("log.user-error", listenCmd.Flags().Lookup("log-user-error")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "log.user-error", err)
	}

	if err := commands.RegisterPostgresOptions(listenCmd); nil != err {
		return nil, err
	}

	return listenCmd, nil
}

func listenRun(_ *cobra.Command, _ []string) {
	addr := net.TCPAddr{
		IP:   net.ParseIP(viper.GetString("server.address")),
		Port: viper.GetInt("server.port"),
		Zone: "",
	}
	errHandler := httperror.HTTPErrorHandler{
		LogUserError: viper.GetBool("log.user-error"),
	}

	sessionManager := scs.New()

	db, err := postgres.New(commands.LoadPostgresConfiguration())
	if err != nil {
		log.Fatal().Err(err).Msg("init postgres connection")
	}

	if err := syncDatabaseSchema(db); nil != err {
		log.Fatal().Err(err).Msg("synchronize database schema")
	}

	userDAO := internalpostgres.UserDAO{Querier: db}
	authenticationMiddleware := internalhttp.AuthenticationMiddleware{
		ErrHandler:     errHandler,
		UserDAO:        userDAO,
		SessionManager: sessionManager,
	}

	r := chi.NewRouter()
	r.Use(authenticationMiddleware.Wrap)
	r.Mount("/", internalhttp.Static)
	r.Route("/api", func(r chi.Router) {
		r.Mount("/auth", api.NewAuth(errHandler, sessionManager).Handler())
		r.Mount("/users", api.Users{ErrHandler: errHandler, UserDAO: userDAO}.Handler())
	})

	if err := internalhttp.Listen(addr, r); nil != err {
		log.Fatal().Msg("listener error")
	}
}

func syncDatabaseSchema(db *sql.DB) error {
	versioner := godatabaseversioner.Versioner{
		Applier: &godatabaseversioner.PostgresVersionApplier{DB: db},
		Listener: godatabaseversioner.ZerologListener{
			Logger: log.Info(),
		},
		Versions: versions.Load(db),
	}
	return versioner.UpgradeToLast()
}
