package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "run",
	Short: "Gochi Template",
	Long:  "Gochi Template",
}

func Execute() error {
	rootCmd.AddCommand(
		restCmd,
	)

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
	os.Exit(0)

	return err
}
