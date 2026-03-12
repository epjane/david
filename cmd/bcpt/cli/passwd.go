package cli

import (
	"fmt"
	"os"
	"syscall"

	"github.com/audstanley/david/app"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/term"
)

var passwdCmd = &cobra.Command{
	Use:   "passwd",
	Short: "Generate a BCrypt password hash",
	Long:  `Generate a BCrypt password hash for use in david configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		password, _ := cmd.Flags().GetString("password")
		cost, _ := cmd.Flags().GetInt("cost")

		if password == "" {
			password = readPassword()
		}

		hash, err := app.GenHashFromPassword(password, cost)
		if err != nil {
			fmt.Printf("Error generating hash: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("Hashed Password: %s\n", hash)
	},
}

func init() {
	RootCmd.AddCommand(passwdCmd)

	passwdCmd.Flags().StringP("password", "p", "", "Password to hash (required)")
	passwdCmd.Flags().IntP("cost", "c", 10, "BCrypt cost factor")
	passwdCmd.MarkFlagRequired("password")

	viper.BindPFlags(passwdCmd.Flags())
}

func readPassword() string {
	fmt.Print("Enter password: ")
	pw, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		fmt.Printf("An error occurred reading the password: %s\n", err)
		os.Exit(1)
	}

	fmt.Println()
	return string(pw)
}
