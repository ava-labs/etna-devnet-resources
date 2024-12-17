package cmd

import (
	"github.com/spf13/cobra"
)

func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use: "manual_etna_evm",
}
