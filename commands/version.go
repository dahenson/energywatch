package commands

import (
	"fmt"
	"github.com/spf13/cobra"
)

const version = `0.1 - Barebones`

var versionCmd = &cobra.Command {
	Use: "version",
	Short: "Display the version.",
	Long: "Display the version number for energywatch",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version + "\n")
	},
}

