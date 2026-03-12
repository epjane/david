package cli

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "david",
	Short: "A simple WebDAV server... extended.",
	Long: `david is a simple WebDAV server that provides:
- Single binary that runs under Windows, Linux and OSX
- Authentication via HTTP-Basic
- CRUD operation permissions
- TLS support
- A simple user management which allows user-directory-jails
- Live config reload
- A CLI tool to generate password hashes (dcrypt - supports bcrypt, argon2, scrypt)`,
	Run: func(cmd *cobra.Command, args []string) {
		RunServer(cmd, args)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		cobra.CheckErr(err)
	}
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func GetServerCmd() *cobra.Command {
	var serverCmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "server" {
			serverCmd = c
			break
		}
	}
	return serverCmd
}
