package data

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func registerResetOption(cmd *cobra.Command) error {
	cmd.Flags().Bool("reset", false, "Drop current database before importing or generating dummy data")
	if err := viper.BindPFlag("sync.reset", cmd.Flags().Lookup("reset")); err != nil {
		return fmt.Errorf("binding parameter %s: %w", "sync.drop-database", err)
	}
	return nil
}
