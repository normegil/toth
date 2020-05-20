package commands

import (
	"fmt"
	"github.com/normegil/postgres"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const ApplicationName = "toth"

func RegisterPostgresOptions(cmd *cobra.Command) error {
	cmd.Flags().String("postgres-address", "localhost", "Postgres server connection URL")
	if err := viper.BindPFlag("postgres.address", cmd.Flags().Lookup("postgres-address")); err != nil {
		return fmt.Errorf("binding parameter %s: %w", "postgres.address", err)
	}

	cmd.Flags().Int("postgres-port", 5432, "Postgres server connection URL")
	if err := viper.BindPFlag("postgres.port", cmd.Flags().Lookup("postgres-port")); err != nil {
		return fmt.Errorf("binding parameter %s: %w", "postgres.port", err)
	}

	cmd.Flags().String("postgres-database", ApplicationName, "Postgres server connection URL")
	if err := viper.BindPFlag("postgres.database", cmd.Flags().Lookup("postgres-database")); err != nil {
		return fmt.Errorf("binding parameter %s: %w", "postgres.database", err)
	}

	cmd.Flags().String("postgres-user", "postgres", "Postgres server connection User")
	if err := viper.BindPFlag("postgres.user", cmd.Flags().Lookup("postgres-user")); err != nil {
		return fmt.Errorf("binding parameter %s: %w", "postgres.user", err)
	}

	cmd.Flags().String("postgres-password", "postgres", "Postgres server connection Password")
	if err := viper.BindPFlag("postgres.password", cmd.Flags().Lookup("postgres-password")); err != nil {
		return fmt.Errorf("binding parameter %s: %w", "postgres.password", err)
	}
	return nil
}

func LoadPostgresConfiguration() postgres.Configuration {
	return postgres.Configuration{
		Address:  viper.GetString("postgres.address"),
		Port:     viper.GetInt("postgres.port"),
		User:     viper.GetString("postgres.user"),
		Password: viper.GetString("postgres.password"),
		Database: viper.GetString("postgres.database"),
		RequiredExtentions: []string{
			"pgcrypto",
		},
	}
}
