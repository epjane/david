package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "bcpt",
	Short: "BCrypt password hash generator",
	Long:  `bcpt is a CLI tool to generate BCrypt password hashes for david configuration.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
