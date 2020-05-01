package main

import (
	"fmt"
	logCfg "github.com/normegil/toth/server/internal/log"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

const ApplicationName = "toth"

var cfgFile string //nolint:gochecknoglobals // Satisfying cobra interface 'OnInitialize' require this global variable

func main() {
	logCfg.Init()
	cobra.OnInitialize(initConfig)

	root, err := root()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not execute command")
	}
	root.PersistentFlags().StringVar(&cfgFile, "config", "", "configuration file")
	if err := root.Execute(); err != nil {
		log.Fatal().Err(err).Msg("Could not execute command")
	}
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cfgDirs := os.Getenv("XDG_CONFIG_DIRS")
		dirs := strings.Split(cfgDirs, ":")
		for _, dir := range dirs {
			viper.AddConfigPath(dir + string(os.PathSeparator) + ApplicationName)
		}

		viper.AddConfigPath("/etc/" + ApplicationName)
		viper.AddConfigPath("$XDG_CONFIG_HOME" + string(os.PathSeparator) + ApplicationName)
		viper.AddConfigPath("$HOME" + string(os.PathSeparator) + "." + ApplicationName)
		viper.AddConfigPath(".")

		viper.SetConfigType("yaml")
		viper.SetConfigName(ApplicationName)
	}

	viper.SetEnvPrefix(strings.ToUpper(ApplicationName) + "_")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, isNotFound := err.(viper.ConfigFileNotFoundError); !isNotFound {
			log.Fatal().Err(err).Msg("could not read configuration")
		}
	}

	logCfg.Configure(viper.GetBool(""))
}

func root() (*cobra.Command, error) {
	rootCmd := &cobra.Command{
		Use:   ApplicationName,
		Short: "Toth is designed to manage a school or a sub-part of a school.",
		Long:  `Toth is designed to manage a school or a sub-part of a school.`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); nil != err {
				log.Fatal().Err(err).Msg("could not print help message")
			}
		},
	}

	rootCmd.PersistentFlags().Bool("color", false, "Colorized & human readable logging")
	if err := viper.BindPFlag("log.color", rootCmd.PersistentFlags().Lookup("color")); err != nil {
		return nil, fmt.Errorf("binding parameter %s: %w", "log.color", err)
	}

	listenCmd, err := listen()
	if err != nil {
		return nil, fmt.Errorf("creating 'listen' command: %w", err)
	}
	rootCmd.AddCommand(listenCmd)

	return rootCmd, nil
}
