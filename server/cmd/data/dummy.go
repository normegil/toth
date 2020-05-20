package data

import (
	"fmt"
	"github.com/normegil/godatabaseversioner"
	"github.com/normegil/postgres"
	"github.com/normegil/toth/server/internal"
	"github.com/normegil/toth/server/internal/commands"
	internalpostgres "github.com/normegil/toth/server/internal/postgres"
	"github.com/normegil/toth/server/internal/postgres/versions"
	"github.com/normegil/toth/server/internal/security"
	"github.com/normegil/toth/server/internal/tools"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func GenerateDummyCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "generate-dummy",
		Short: "Generate dummy data",
		Long:  `Generate dummy data`,
		Args:  cobra.NoArgs,
		Run:   generateDummyData,
	}

	if err := commands.RegisterPostgresOptions(cmd); nil != err {
		return nil, fmt.Errorf("registering postgres options: %w", err)
	}

	if err := registerResetOption(cmd); err != nil {
		return nil, fmt.Errorf("registering reset options: %w", err)
	}

	return cmd, nil
}

func generateDummyData(_ *cobra.Command, _ []string) {
	pgCfg := commands.LoadPostgresConfiguration()

	dropDatabase := viper.GetBool("sync.reset")
	if dropDatabase {
		if err := postgres.DropDatabase(pgCfg); nil != err {
			log.Fatal().Err(err).Msg("Could not drop database")
		}
	}

	db, err := postgres.New(pgCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Initialize database connection")
	}
	defer tools.Close(db)

	applier := godatabaseversioner.PostgresVersionApplier{DB: db}
	versioner := godatabaseversioner.NewVersioner(applier, versions.Load(db))

	if err = versioner.UpgradeToLast(); nil != err {
		log.Fatal().Err(err).Int("version", versioner.LastVersion()).Msg("Sync database version")
	}

	mobarvauxUser, err := security.NewUser("Marie-Odile Barvaux", "mobarvaux@gmail.com", "test")
	if err != nil {
		log.Fatal().Err(err).Msg("generating user")
	}

	users := []internal.User{
		*mobarvauxUser,
	}
	userDAO := internalpostgres.UserDAO{Querier: db}
	for _, user := range users {
		if err := userDAO.Insert(user); nil != err {
			log.Fatal().Err(err).Str("user", user.Name).Msg("inserting user")
		}
	}
}
