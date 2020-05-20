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
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ExportCmd() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:       "export",
		Short:     "Export data to specified file",
		Long:      `Export data to specified file`,
		ValidArgs: []string{"exported-file-name.yml"},
		Args:      cobra.ExactArgs(1),
		Run:       exportData,
	}

	if err := commands.RegisterPostgresOptions(cmd); nil != err {
		return nil, fmt.Errorf("registering postgres options: %w", err)
	}

	return cmd, nil
}

func exportData(_ *cobra.Command, args []string) {
	targetFilePath := args[0]

	pgCfg := commands.LoadPostgresConfiguration()

	db, err := postgres.New(pgCfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Initialize database connection")
	}
	defer tools.Close(db)

	applier := godatabaseversioner.PostgresVersionApplier{DB: db}
	versioner := godatabaseversioner.NewVersioner(applier, versions.Load(db))
	if err = versioner.UpgradeToLast(); nil != err {
		log.Fatal().Err(err).Msg("upgrading database")
	}

	currentVersion, err := versioner.CurrentVersion()
	if err != nil {
		log.Fatal().Err(err).Msg("loading current database version")
	}

	dao := internalpostgres.UserDAO{Querier: db}
	users, err := dao.LoadAll()
	if err != nil {
		log.Fatal().Err(err).Msg("load all users")
	}

	model := syncPkg.Model{
		Version: currentVersion,
		Data: syncPkg.Data{
			Users: users,
		},
	}
	marshalled, err := yaml.Marshal(model)
	if err != nil {
		log.Fatal().Err(err).Msgf("Could not marshal %+v", model)
	}

	if err = ioutil.WriteFile(targetFilePath, marshalled, 0777); nil != err {
		log.Fatal().Err(err).Msgf("Could not write test file")
	}
}
