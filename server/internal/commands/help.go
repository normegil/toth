package commands

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func PrintHelp(cmd *cobra.Command, _ []string) {
	if err := cmd.Help(); nil != err {
		log.Fatal().Err(err).Msg("could not print help message")
	}
}
