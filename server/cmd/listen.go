package main

import (
	"database/sql"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/normegil/godatabaseversioner"
	"github.com/normegil/postgres"
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

	listenCmd.Flags().String("postgres-address", "localhost", "Postgres server connection URL")
	if err := viper.BindPFlag("postgres.address", listenCmd.Flags().Lookup("postgres-address")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "postgres.address", err)
	}

	listenCmd.Flags().Int("postgres-port", 5432, "Postgres server connection URL")
	if err := viper.BindPFlag("postgres.port", listenCmd.Flags().Lookup("postgres-port")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "postgres.port", err)
	}

	listenCmd.Flags().String("postgres-database", ApplicationName, "Postgres server connection URL")
	if err := viper.BindPFlag("postgres.database", listenCmd.Flags().Lookup("postgres-database")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "postgres.database", err)
	}

	listenCmd.Flags().String("postgres-user", "postgres", "Postgres server connection User")
	if err := viper.BindPFlag("postgres.user", listenCmd.Flags().Lookup("postgres-user")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "postgres.user", err)
	}

	listenCmd.Flags().String("postgres-password", "postgres", "Postgres server connection Password")
	if err := viper.BindPFlag("postgres.password", listenCmd.Flags().Lookup("postgres-password")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "postgres.password", err)
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

	db, err := postgres.New(postgres.Configuration{
		Address:  viper.GetString("postgres.address"),
		Port:     viper.GetInt("postgres.port"),
		User:     viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
		Database: viper.GetString("postgres.database"),
		RequiredExtentions: []string{
			"pgcrypto",
		},
	})
	if err != nil {
		log.Fatal().Err(err).Msg("init postgres connection")
	}

	if err := syncDatabaseSchema(db); nil != err {
		log.Fatal().Err(err).Msg("synchronize database schema")
	}

	authenticationMiddleware := internalhttp.AuthenticationMiddleware{
		ErrHandler:     errHandler,
		UserDAO:        internalpostgres.UserDAO{Querier: db},
		SessionManager: sessionManager,
	}

	r := chi.NewRouter()
	r.Use(authenticationMiddleware.Wrap)
	r.Mount("/", internalhttp.Static)
	r.Route("/api", func(r chi.Router) {
		r.Mount("/auth", api.NewAuth(errHandler, sessionManager).Handler())
		r.Mount("/users", api.Users{ErrHandler: errHandler}.Handler())
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
