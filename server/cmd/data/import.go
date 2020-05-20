package data

import (
	"fmt"
	"github.com/normegil/godatabaseversioner"
	"github.com/normegil/postgres"
	"github.com/normegil/toth/server/internal/commands"
	internalpostgres "github.com/normegil/toth/server/internal/postgres"
	"github.com/normegil/toth/server/internal/postgres/versions"
	syncPkg "github.com/normegil/toth/server/internal/sync"
	"github.com/normegil/toth/server/internal/tools"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ImportCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import data to storages",
		Long:  `Import data to storages`,
		Args:  cobra.ExactArgs(1),
		Run:   importData,
	}

	if err := commands.RegisterPostgresOptions(cmd); nil != err {
		return nil, fmt.Errorf("registering postgres options: %w", err)
	}

	if err := registerResetOption(cmd); err != nil {
		return nil, fmt.Errorf("registering reset options: %w", err)
	}

	return cmd, nil
}

func importData(_ *cobra.Command, args []string) {
	importedFilePath := args[0]

	inputContent, err := ioutil.ReadFile(importedFilePath)
	if err != nil {
		log.Fatal().Err(err).Str("path", importedFilePath).Msg("reading input file")
	}
	var model syncPkg.Model
	if err = yaml.Unmarshal(inputContent, &model); nil != err {
		log.Fatal().Err(err).Str("path", importedFilePath).Msg("unmarshal input file")
	}

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

	if err = versioner.Sync(model.Version); nil != err {
		log.Fatal().Err(err).Int("version", model.Version).Msg("Sync database version")
	}

	userDAO := internalpostgres.UserDAO{Querier: db}
	for _, user := range model.Data.Users {
		if err := userDAO.Insert(user); nil != err {
			log.Fatal().Err(err).Interface("user", user).Msg("inserting user")
		}
	}
}
