package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "enshamir",
	Short: "enshamir is a tool for encrypting and splitting secrets into shares",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

const saltFileName = "MUST-BACK-UP-SALT"
const shareFilePrefix = "SPLITTED-SECRET"

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
