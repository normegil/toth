package main

import (
	"fmt"
	dataCmdPkg "github.com/normegil/toth/server/cmd/data"
	"github.com/normegil/toth/server/internal/commands"
	"github.com/spf13/cobra"
)

func data() (*cobra.Command, error) {
	dataCmd := &cobra.Command{
		Use:   "data",
		Short: "Manage data used by toth",
		Long:  `Manage data used by toth`,
		Run:   commands.PrintHelp,
	}
	if err := commands.RegisterPostgresOptions(dataCmd); nil != err {
		return nil, err
	}

	importCmd, err := dataCmdPkg.ImportCmd()
	if err != nil {
		return nil, fmt.Errorf("creating 'import' command: %w", err)
	}
	dataCmd.AddCommand(importCmd)

	exportCmd, err := dataCmdPkg.ExportCmd()
	if err != nil {
		return nil, fmt.Errorf("creating 'export' command: %w", err)
	}
	dataCmd.AddCommand(exportCmd)

	generateDummyCmd, err := dataCmdPkg.GenerateDummyCmd()
	if err != nil {
		return nil, fmt.Errorf("creating 'generate-dummy' command: %w", err)
	}
	dataCmd.AddCommand(generateDummyCmd)

	return dataCmd, nil
}
