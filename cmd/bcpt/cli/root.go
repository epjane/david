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

func GetRootCmd() *cobra.Command {
	return RootCmd
}

func GetPasswdCmd() *cobra.Command {
	var passwdCmd *cobra.Command
	for _, c := range RootCmd.Commands() {
		if c.Use == "passwd" {
			passwdCmd = c
			break
		}
	}
	return passwdCmd
}
